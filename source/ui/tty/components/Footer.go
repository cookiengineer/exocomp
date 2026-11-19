package components

import "exocomp/ui/tty/styles"
import "exocomp/utils/chat/commands"
import utils_tty "exocomp/utils/tty"
import "strings"
import "time"

type Footer struct {
	Catalog     *commands.Catalog
	Textarea    *Textarea
	Palette     *Palette
	Status      string
	Role        string
	Focused     bool
	cursorX     int
	cursorY     int
	cursorShown bool
}

func NewFooter(catalog *commands.Catalog) *Footer {

	return &Footer{
		Catalog:  catalog,
		Textarea: NewTextarea(),
		Palette:  NewPalette(),
		Role:     "user",
		Focused:  true,
	}

}

func (footer *Footer) AcceptSuggestion() {

	item, ok := footer.Palette.GetSelectedItem()

	if ok == false {
		return
	}

	if item.Kind == "tool" {

		footer.Textarea.SetText("/" + item.Value + " ")
		footer.RefreshSuggestions()

	} else if item.Kind == "parameter" {

		text := footer.Textarea.GetText()

		if strings.HasSuffix(text, " ") == true || text == "" {
			footer.Textarea.Insert(item.Value + "=")
		} else {
			footer.Textarea.Insert(" " + item.Value + "=")
		}

		footer.RefreshSuggestions()

	} else if item.Kind == "command" {

		footer.Textarea.SetText(item.Value + " ")
		footer.RefreshSuggestions()

	}

}

func (footer *Footer) SetCatalog(catalog *commands.Catalog) {

	footer.Catalog = catalog
	footer.RefreshSuggestions()

}

func (footer *Footer) GetText() string {
	return footer.Textarea.GetText()
}

func (footer *Footer) Clear() {

	footer.Textarea.Clear()
	footer.Palette.Hide()

}

func (footer *Footer) Width(width int) int {
	return width
}

func (footer *Footer) Height(width int) int {

	textarea_width := max(1, width-5)
	lines, _, _ := footer.Textarea.Wrap(textarea_width)
	textarea_height := min(10, max(1, len(lines)))
	palette_height := 0

	if footer.Palette.Visible == true {
		palette_height = footer.Palette.Height(textarea_width)
	}

	return 1 + textarea_height + palette_height + 1 + 1

}

func (footer *Footer) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

	if rect.Width <= 0 || rect.Height <= 0 {
		return
	}

	status_y := rect.Y + rect.Height - 1
	cap_y := rect.Y + rect.Height - 2
	panel_height := max(0, rect.Height-2)
	border_style := styles.PromptBorder

	if footer.Focused == true {
		border_style = styles.PromptBorderFocus
	}

	if panel_height > 0 {

		panel_rect := utils_tty.Rect{
			X:      rect.X,
			Y:      rect.Y,
			Width:  rect.Width,
			Height: panel_height,
		}

		buffer.Fill(panel_rect, ' ', styles.PromptPanel)

		for row := rect.Y; row < rect.Y+panel_height; row++ {
			buffer.Set(rect.X, row, "┃", border_style)
		}

	}

	buffer.Fill(utils_tty.Rect{
		X:      rect.X,
		Y:      cap_y,
		Width:  rect.Width,
		Height: 1,
	}, ' ', styles.Background)

	buffer.Set(rect.X, cap_y, "╹", border_style)

	for x := rect.X + 1; x < rect.X+rect.Width; x++ {
		buffer.Set(x, cap_y, "▀", styles.PromptCap)
	}

	buffer.Fill(utils_tty.Rect{
		X:      rect.X,
		Y:      status_y,
		Width:  rect.Width,
		Height: 1,
	}, ' ', styles.Background)

	textarea_width := max(1, rect.Width-5)
	lines, cursor_line, cursor_column := footer.Textarea.Wrap(textarea_width)
	textarea_desired := min(10, max(1, len(lines)))
	content_height := max(0, panel_height-1)
	palette_desired := 0

	if footer.Palette.Visible == true {
		palette_desired = footer.Palette.Height(textarea_width)
	}

	palette_height := min(palette_desired, max(0, content_height-1))
	textarea_height := 0

	if content_height > 0 {
		textarea_height = min(textarea_desired, max(1, content_height-palette_height))
	}

	if textarea_height > 0 {

		if cursor_line < footer.Textarea.scroll {
			footer.Textarea.scroll = cursor_line
		} else if cursor_line >= footer.Textarea.scroll+textarea_height {
			footer.Textarea.scroll = cursor_line - textarea_height + 1
		}

	}

	textarea_rect := utils_tty.Rect{
		X:      rect.X + 3,
		Y:      rect.Y + 1,
		Width:  textarea_width,
		Height: textarea_height,
	}

	if textarea_height > 0 && footer.Textarea.GetText() == "" {

		placeholder := "Ask anything, / for commands and tools, @ for context…"
		buffer.WriteString(textarea_rect.X, textarea_rect.Y, utils_tty.Truncate(placeholder, textarea_rect.Width), styles.PromptMuted)

	} else if textarea_height > 0 {

		for row := 0; row < textarea_rect.Height; row++ {

			index := footer.Textarea.scroll + row

			if index >= len(lines) {
				break
			}

			line := lines[index]
			buffer.WriteString(textarea_rect.X, textarea_rect.Y+row, utils_tty.Truncate(line.Text, textarea_rect.Width), styles.PromptText)

		}

	}

	if footer.Palette.Visible == true && palette_height > 0 {

		palette_rect := utils_tty.Rect{
			X:      rect.X + 3,
			Y:      textarea_rect.Y + textarea_height,
			Width:  textarea_width,
			Height: palette_height,
		}

		footer.Palette.Render(buffer, palette_rect, elapsed)

	}

	status := sanitizeFooterStatus(footer.Status)
	left, busy := footerActivity(status)
	bindings := footerBindings(status)

	if busy == true {
		left = spinnerFrame(elapsed) + " " + left
	}

	right_width := footerBindingsWidth(bindings)
	left_width := max(0, rect.Width-5-right_width-2)
	left = utils_tty.TruncateRight(left, left_width)

	buffer.WriteString(rect.X+3, status_y, left, styles.TextMuted)
	writeRightBindings(buffer, rect.X+rect.Width-2, status_y, rect.X+3, bindings)

	footer.cursorShown = footer.Focused && textarea_height > 0
	footer.cursorX = textarea_rect.X + cursor_column
	footer.cursorY = textarea_rect.Y + (cursor_line - footer.Textarea.scroll)

	if textarea_height > 0 {

		if footer.cursorX < textarea_rect.X {
			footer.cursorX = textarea_rect.X
		}

		if footer.cursorX >= textarea_rect.X+textarea_rect.Width {
			footer.cursorX = textarea_rect.X + textarea_rect.Width - 1
		}

		if footer.cursorY < textarea_rect.Y {
			footer.cursorY = textarea_rect.Y
		}

		if footer.cursorY >= textarea_rect.Y+textarea_rect.Height {
			footer.cursorY = textarea_rect.Y + textarea_rect.Height - 1
		}

	} else {

		footer.cursorX = textarea_rect.X
		footer.cursorY = textarea_rect.Y

	}

}

func (footer *Footer) GetCursor() (int, int, bool) {
	return footer.cursorX, footer.cursorY, footer.cursorShown
}

func (footer *Footer) HandleKey(key utils_tty.Key) bool {

	if key.Type == utils_tty.KeyPaste {

		footer.Textarea.Insert(key.Text)
		footer.RefreshSuggestions()

		return false

	}

	if key.Type == utils_tty.KeyRune && key.Ctrl == true && key.Rune == 's' {
		return true
	}

	if footer.Palette.Visible == true {

		switch key.Type {

		case utils_tty.KeyArrowUp:
			footer.Palette.Move(-1)
			return false
		case utils_tty.KeyArrowDown:
			footer.Palette.Move(1)
			return false
		case utils_tty.KeyTab:
			footer.AcceptSuggestion()
			return false
		case utils_tty.KeyEscape:
			footer.Palette.Hide()
			return false
		case utils_tty.KeyEnter:
			footer.AcceptSuggestion()
			return false

		}

	}

	switch key.Type {

	case utils_tty.KeyEnter:

		footer.Textarea.InsertNewline()
		footer.RefreshSuggestions()

		return false

	case utils_tty.KeyTab:

		if footer.Palette.Visible == true {
			footer.AcceptSuggestion()
		} else {

			footer.Textarea.Insert("    ")
			footer.RefreshSuggestions()

		}

	case utils_tty.KeyBackspace:

		if footer.Textarea.Backspace() == true {
			footer.RefreshSuggestions()
		}

	case utils_tty.KeyDelete:
		footer.Textarea.Delete()
		footer.RefreshSuggestions()

	case utils_tty.KeyArrowLeft:
		footer.Textarea.MoveLeft()

	case utils_tty.KeyArrowRight:
		footer.Textarea.MoveRight()

	case utils_tty.KeyArrowUp:
		footer.Textarea.MoveUp()

	case utils_tty.KeyArrowDown:
		footer.Textarea.MoveDown()

	case utils_tty.KeyHome:
		footer.Textarea.Home()

	case utils_tty.KeyEnd:
		footer.Textarea.End()

	case utils_tty.KeyEscape:
		footer.Palette.Hide()

	case utils_tty.KeyRune:

		if key.Ctrl == true {

			if key.Rune == 'j' {

				footer.Textarea.InsertNewline()
				footer.RefreshSuggestions()

			}

		} else {

			footer.Textarea.Insert(string(key.Rune))
			footer.RefreshSuggestions()

		}

	}

	return false

}

func (footer *Footer) RefreshSuggestions() {

	text := footer.Textarea.GetText()

	if footer.Catalog == nil || strings.HasPrefix(text, "/") == false {

		footer.Palette.Hide()
		return

	}

	if strings.ContainsAny(strings.TrimPrefix(text, "/"), " \t\n") == true {

		items := footer.Catalog.Complete(text, 8)
		footer.Palette.Set(items)

		return

	}

	items := footer.Catalog.Complete(text, 12)
	result := make([]commands.Suggestion, 0)
	result = append(result, items...)

	prefix := strings.ToLower(strings.TrimPrefix(text, "/"))

	for _, builtin := range Builtins {

		name := strings.TrimPrefix(builtin.Name, "/")

		if prefix == "" || strings.HasPrefix(strings.ToLower(name), prefix) == true {
			result = append(result, commands.Suggestion{
				Kind:        "command",
				Value:       builtin.Name,
				Description: builtin.Description,
			})
		}

	}

	footer.Palette.Set(result)

}
