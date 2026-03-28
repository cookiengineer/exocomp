package ollama

import "golang.org/x/term"
import "exocomp/schemas"
import "fmt"
import "os"
import "strings"
import "sync"
import "time"

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

type Renderer struct {
	Prompt    string
	Session   *Session
	mutex     *sync.Mutex
	rendered  int
	old_state *term.State
}

func NewRenderer(session *Session) *Renderer {

	stdin_num := int(os.Stdin.Fd())
	old_state, err := term.MakeRaw(stdin_num)

	if err == nil {

		return &Renderer{
			Prompt:    "",
			Session:   session,
			mutex:     &sync.Mutex{},
			rendered:  0,
			old_state: old_state,
		}

	} else {

		return &Renderer{
			Prompt:    "",
			Session:   session,
			mutex:     &sync.Mutex{},
			rendered:  0,
			old_state: nil,
		}

	}

}

func (renderer *Renderer) Destroy() {

	if renderer.old_state != nil {
		term.Restore(int(os.Stdin.Fd()), renderer.old_state)
	}

}

func (renderer *Renderer) InputLoop() {

	buffer := make([]byte, 1)

	for {

		length, err0 := os.Stdin.Read(buffer)

		if err0 == nil && length > 0 {

			chr := buffer[0]

			renderer.mutex.Lock()

			switch chr {
			case 3: // Ctrl+C

				// Stop InputLoop
				return

			case 13: // Enter

				user_prompt := strings.TrimSpace(renderer.Prompt)

				if user_prompt == "bye!" || user_prompt == "quit" || user_prompt == "exit" || user_prompt == "exit()" {

					renderer.Prompt = ""
					return

				} else {

					go renderer.Session.Query(schemas.Message{
						Role:    "user",
						Content: renderer.Prompt,
					})
					renderer.Prompt = ""

				}

			case 127: // Backspace

				if len(renderer.Prompt) > 0 {
					renderer.Prompt = renderer.Prompt[0:len(renderer.Prompt)-1]
				}

			default:

				renderer.Prompt += string(chr)

			}

			renderer.mutex.Unlock()

		} else {
			continue
		}

	}

}

func (renderer *Renderer) RenderLoop() {

	// 10 FPS should be enough
	ticker := time.NewTicker(1000 / 10 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {

		renderer.mutex.Lock()

		renderer.ClearScreen()
		renderer.RenderMessages(renderer.Session.Messages)
		renderer.RenderPrompt()

		renderer.mutex.Unlock()

	}

}

func (renderer *Renderer) ClearScreen() {

	fmt.Fprint(os.Stdout, "\033[H\033[2J")
	os.Stdout.Sync()

}

func (renderer *Renderer) RenderMessages(messages []*schemas.Message) {

	for _, message := range messages {

		color   := ColorReset
		role    := message.Role
		content := formatContent(message.Content)

		switch message.Role {
		case "user":
			color = ColorGreen
		case "assistant":
			color = ColorBlue
		case "tool":
			color = ColorRed
		case "system":

			if renderer.Session.Config.Verbose == true {
				color = ColorYellow
			} else {
				color = ColorReset
			}

		default:
			color = ColorReset
		}

		if color != ColorReset && len(content) > 0 {

			if len(content) > 1 {

				fmt.Fprintf(os.Stdout, "\r%s[%s]%s:\n", color, role, ColorReset)

				for _, line := range content {
					fmt.Fprintf(os.Stdout, "\r%s|%s %s\n", color, ColorReset, line)
				}

				os.Stdout.Sync()

			} else {

				fmt.Fprintf(os.Stdout, "\r%s[%s]%s: %s\n", color, role, ColorReset, content[0])
				os.Stdout.Sync()

			}

		}

	}

}

func (renderer *Renderer) RenderPrompt() {

	model  := "unknown"
	prompt := renderer.Prompt

	if renderer.Session != nil && renderer.Session.Config != nil {
		model = renderer.Session.Config.Model
	}

	fmt.Fprintf(os.Stdout, "\r%s[to %s]%s > %s", ColorGreen, model, ColorReset, prompt)
	os.Stdout.Sync()

}
