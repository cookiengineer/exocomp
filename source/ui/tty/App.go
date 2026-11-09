package tty

import "exocomp/engine"
import "exocomp/schemas"
import "exocomp/tools"
import "exocomp/types"
import "exocomp/ui/tty/components"
import "exocomp/ui/tty/styles"
import "exocomp/utils/chat/commands"
import utils_tty "exocomp/utils/tty"
import "os"
import "os/signal"
import "strings"
import "sync"
import "syscall"
import "time"

type App struct {
	Session     *engine.Session
	Terminal    *utils_tty.Terminal
	Input       *utils_tty.Input
	Screen      *utils_tty.Screen
	Header      *components.Header
	Transcript  *components.Transcript
	Sidebar     *components.Sidebar
	Footer      *components.Footer
	Human       components.Dialog
	Role        string
	status      string
	statusMutex *sync.Mutex
	started     time.Time
	lastHuman   time.Time
	catalog     *commands.Catalog
	quit        chan struct{}
	quitOnce    *sync.Once
}

func NewApp(session *engine.Session) *App {

	catalog := commands.NewCatalog(session.GetToolSchemas())

	return &App{
		Session:     session,
		Terminal:    utils_tty.NewTerminal(),
		Header:      components.NewHeader(session),
		Transcript:  components.NewTranscript(session),
		Sidebar:     components.NewSidebar(session),
		Footer:      components.NewFooter(catalog),
		Role:        "user",
		statusMutex: &sync.Mutex{},
		started:     time.Now(),
		catalog:     catalog,
		quit:        make(chan struct{}),
		quitOnce:    &sync.Once{},
	}

}

func (app *App) Init() {

	_ = app.Run()

}

func (app *App) SetRole(role string) {

	if role == "user" || role == "assistant" {

		app.Role = role
		app.Header.Role = role
		app.Footer.Role = role

	}

}

func (app *App) Destroy() {

	app.stop()
	app.Terminal.Leave()

}

func (app *App) Run() error {

	err0 := app.Terminal.Enter()

	if err0 == nil {

		defer app.Terminal.Leave()

		width, height := app.Terminal.Size()
		app.Screen = utils_tty.NewScreen(app.Terminal, width, height)
		app.Input = utils_tty.NewInput(app.Terminal)

		defer app.Input.Close()

		keys := make(chan utils_tty.Key, 32)

		go func() {

			for {

				key, err := app.Input.Read()

				if err != nil {

					select {

					case keys <- utils_tty.Key{Type: utils_tty.KeyEOF}:
					default:

					}

					return

				}

				keys <- key

			}

		}()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(signals)

		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		changes := app.Session.Changes()

		app.render(time.Now())

		for {

			select {

			case key := <-keys:
				app.handleKey(key)
				app.render(time.Now())

			case <-changes:
				app.refresh()
				app.render(time.Now())

			case now := <-ticker.C:
				app.tick(now)
				app.render(now)

			case <-signals:
				app.stop()
				return nil

			case <-app.quit:
				return nil

			}

		}

	} else {
		return err0
	}

}

func (app *App) stop() {

	app.quitOnce.Do(func() {
		close(app.quit)
	})

}

func (app *App) refresh() {

	app.catalog = commands.NewCatalog(app.Session.GetToolSchemas())
	app.Footer.SetCatalog(app.catalog)
	app.Sidebar.Refresh(time.Now(), true)
	app.Transcript.Invalidate()

}

func (app *App) tick(now time.Time) {

	if app.Screen != nil {

		width, height := app.Terminal.Size()

		if width != app.Screen.Width() || height != app.Screen.Height() {
			app.Screen.Resize(width, height)
		}

	}

	app.Sidebar.Refresh(now, false)

	if app.Human == nil && now.Sub(app.lastHuman) > 300*time.Millisecond {

		app.lastHuman = now
		app.checkHuman()

	}

}

func (app *App) checkHuman() {

	tool := app.Session.GetTool("humans")

	if tool != nil {

		humans_tool, ok := tool.(*tools.Humans)

		if ok == true {

			pending := humans_tool.Snapshot()

			if len(pending) > 0 {

				if pending[0].Type == "Choose" {
					app.Human = components.NewChoiceDialog(pending[0])
				} else {
					app.Human = components.NewQuestionDialog(pending[0])
				}

				app.Footer.Focused = false

			}

		}

	}

}

func (app *App) handleKey(key utils_tty.Key) {

	if key.Type == utils_tty.KeyEOF {

		app.stop()
		return

	}

	if app.Human != nil {

		answer, submit, cancel := app.Human.HandleKey(key)

		if submit == true {

			question := app.Human.GetQuestion()
			app.Human = nil
			app.Footer.Focused = true
			app.setStatus("Answering human question...")

			go func() {

				err := app.Session.CallTool("", "humans", "Answer", map[string]any{
					"question": question,
					"answer":   answer,
				})

				if err == nil {
					app.setStatus("")
				} else {
					app.setStatus(err.Error())
				}

			}()

		} else if cancel == true {

			app.Human = nil
			app.Footer.Focused = true

		}

		return

	}

	if key.Type == utils_tty.KeyRune && key.Ctrl == true {

		switch key.Rune {

		case 'q':
			app.stop()
			return

		case 'l':
			app.Screen.Resize(app.Screen.Width(), app.Screen.Height())
			return

		case 'c':

			if strings.TrimSpace(app.Footer.GetText()) != "" {

				app.Footer.Clear()
				return

			}

			app.setStatus("Press Ctrl+Q to quit.")
			return

		case 'y':
			_ = app.Terminal.Copy(app.Footer.GetText())
			app.setStatus("Copied composer through OSC52.")
			return

		}

	}

	switch key.Type {

	case utils_tty.KeyArrowUp:

		if key.Ctrl == true || key.Alt == true {

			app.Transcript.ScrollBy(1)
			return

		}

	case utils_tty.KeyArrowDown:

		if key.Ctrl == true || key.Alt == true {

			app.Transcript.ScrollBy(-1)
			return

		}

	case utils_tty.KeyHome:

		if key.Ctrl == true {

			app.Transcript.ScrollToTop()
			return

		}

	case utils_tty.KeyEnd:

		if key.Ctrl == true {

			app.Transcript.ScrollToBottom()
			return

		}

	case utils_tty.KeyPageUp:
		app.Transcript.ScrollBy(max(5, app.Screen.Height()/2))
		return

	case utils_tty.KeyPageDown:
		app.Transcript.ScrollBy(-max(5, app.Screen.Height()/2))
		return

	}

	if app.Footer.HandleKey(key) == true {

		text := app.Footer.GetText()

		if strings.TrimSpace(text) == "" {
			return
		}

		app.submit(text)

	}

}

func (app *App) submit(text string) {

	trimmed := strings.TrimSpace(text)

	if trimmed == "" {
		return
	}

	if strings.HasPrefix(trimmed, "/") == true {

		command, err0 := commands.Resolve(app.catalog, trimmed)

		if err0 == nil {

			app.Footer.Clear()
			app.setStatus("Running " + command.Name + "...")

			go func() {

				err := app.Session.CallTool("", command.Name, command.Method, command.Arguments)

				if err == nil {
					app.setStatus("")
				} else {
					app.setStatus(err.Error())
				}

			}()

			return

		}

		if app.handleBuiltin(trimmed) == true {
			return
		}

		app.setStatus(err0.Error())
		return

	}

	if app.Session.IsWaiting() == true {

		app.setStatus("Session is busy, waiting for the LLM response.")
		return

	}

	app.Footer.Clear()
	app.setStatus("Waiting for the LLM response...")

	go func() {

		err := app.Session.SendChatRequest(schemas.Message{
			Role:    "user",
			Content: trimmed,
		})

		if err == nil {
			app.setStatus("")
		} else {
			app.setStatus(err.Error())
		}

	}()

}

func (app *App) handleBuiltin(text string) bool {

	parts := strings.Fields(text)

	if len(parts) == 0 {
		return false
	}

	switch parts[0] {

	case "/help":
		app.Footer.Clear()
		names := make([]string, 0)

		for _, builtin := range components.Builtins {
			names = append(names, builtin.Name)
		}

		app.setStatus("Commands: " + strings.Join(names, " "))
		return true

	case "/clear":
		app.Footer.Clear()
		app.Transcript.Clear()
		app.setStatus("Cleared the transcript view.")
		return true

	case "/copy":
		_ = app.Terminal.Copy(app.Footer.GetText())
		app.Footer.Clear()
		app.setStatus("Copied composer through OSC52.")
		return true

	case "/sidebar":
		app.Footer.Clear()
		app.Sidebar.Toggle()
		return true

	case "/quit":
		app.stop()
		return true

	default:
		return false

	}

}

func (app *App) setStatus(value string) {

	app.statusMutex.Lock()
	app.status = value
	app.statusMutex.Unlock()

}

func (app *App) statusValue() string {

	app.statusMutex.Lock()
	value := app.status
	app.statusMutex.Unlock()

	return value

}

func (app *App) render(now time.Time) {

	if app.Screen == nil {
		return
	}

	app.Footer.Status = app.statusValue()

	width := app.Screen.Width()
	height := app.Screen.Height()
	buffer := app.Screen.Buffer()
	buffer.Clear(styles.Background)

	header_height := app.Header.Height(width)
	header_rect := utils_tty.Rect{X: 0, Y: 0, Width: width, Height: min(header_height, height)}
	app.Header.Render(buffer, header_rect, now.Sub(app.started))

	sidebar_width := 0

	if app.Sidebar.IsVisible(width) == true {
		sidebar_width = app.Sidebar.Width(width)
	}

	main_width := width - sidebar_width
	footer_height := min(app.Footer.Height(main_width), max(1, height-header_rect.Height))
	content_height := max(0, height-header_rect.Height-footer_height)

	transcript_rect := utils_tty.Rect{
		X:      0,
		Y:      header_rect.Height,
		Width:  main_width,
		Height: content_height,
	}
	app.Transcript.Render(buffer, transcript_rect, now.Sub(app.started))

	footer_rect := utils_tty.Rect{
		X:      0,
		Y:      header_rect.Height + content_height,
		Width:  main_width,
		Height: footer_height,
	}
	app.Footer.Render(buffer, footer_rect, now.Sub(app.started))

	if sidebar_width > 0 {

		sidebar_rect := utils_tty.Rect{
			X:      main_width,
			Y:      header_rect.Height,
			Width:  sidebar_width,
			Height: max(0, height-header_rect.Height),
		}
		app.Sidebar.Render(buffer, sidebar_rect, now.Sub(app.started))

	}

	if app.Human != nil {

		dialog_width := app.Human.Width(width)
		dialog_height := min(app.Human.Height(dialog_width), max(4, height-2))
		dialog_rect := utils_tty.Rect{
			X:      max(0, (width-dialog_width)/2),
			Y:      max(0, (height-dialog_height)/2),
			Width:  dialog_width,
			Height: dialog_height,
		}

		app.Human.Render(buffer, dialog_rect, now.Sub(app.started))

		cursor_x, cursor_y, cursor_visible := app.Human.GetCursor()

		if cursor_visible == true {
			app.Screen.SetCursor(cursor_x, cursor_y, true)
		} else {
			app.Screen.SetCursor(0, 0, false)
		}

	} else {

		cursor_x, cursor_y, cursor_visible := app.Footer.GetCursor()

		if cursor_visible == true {
			app.Screen.SetCursor(cursor_x, cursor_y, true)
		} else {
			app.Screen.SetCursor(0, 0, false)
		}

	}

	_ = app.Screen.Flush()

}

func (app *App) SessionAgent() *types.Agent {

	if app.Session != nil {
		return app.Session.Agent
	} else {
		return nil
	}

}
