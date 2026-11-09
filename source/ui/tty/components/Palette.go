package components

import "exocomp/ui/tty/styles"
import "exocomp/utils/chat/commands"
import utils_tty "exocomp/utils/tty"
import "fmt"
import "time"

type Palette struct {
	Visible  bool
	Items    []commands.Suggestion
	Selected int
}

func NewPalette() *Palette {

	return &Palette{
		Items: make([]commands.Suggestion, 0),
	}

}

func (palette *Palette) Set(items []commands.Suggestion) {

	palette.Items = items
	palette.Visible = len(items) > 0

	if palette.Selected >= len(items) {
		palette.Selected = max(0, len(items)-1)
	}

}

func (palette *Palette) Hide() {

	palette.Visible = false
	palette.Selected = 0

}

func (palette *Palette) Move(delta int) {

	if len(palette.Items) == 0 {
		return
	}

	palette.Selected = (palette.Selected + delta + len(palette.Items)) % len(palette.Items)

}

func (palette *Palette) GetSelectedItem() (commands.Suggestion, bool) {

	if palette.Selected >= 0 && palette.Selected < len(palette.Items) {
		return palette.Items[palette.Selected], true
	} else {
		return commands.Suggestion{}, false
	}

}

func (palette *Palette) Width(width int) int {
	return width
}

func (palette *Palette) Height(width int) int {

	if palette.Visible == false {
		return 0
	} else {
		return min(8, len(palette.Items))
	}

}

func (palette *Palette) Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration) {

	if palette.Visible == false || rect.Height <= 0 || rect.Width <= 0 {
		return
	}

	buffer.Fill(rect, ' ', styles.Palette)

	end := min(len(palette.Items), rect.Height)

	for row := 0; row < end; row++ {

		item := palette.Items[row]
		selected := row == palette.Selected

		style := styles.Palette

		if selected == true {
			style = styles.Selection
		}

		label := item.Value

		if item.Kind == "parameter" {

			if item.Required == true {
				label = fmt.Sprintf("%s=*", item.Value)
			} else {
				label = fmt.Sprintf("%s=", item.Value)
			}

		}

		description := item.Description

		if item.Type != "" {
			description = fmt.Sprintf("[%s] %s", item.Type, description)
		}

		width := rect.Width - 2
		text := utils_tty.PadRight(label, min(max(12, width/3), width))

		if description != "" && width > utils_tty.StringWidth(label)+4 {
			text = utils_tty.Truncate(text+" "+description, width)
		} else {
			text = utils_tty.Truncate(text, width)
		}

		buffer.Fill(utils_tty.Rect{X: rect.X, Y: rect.Y + row, Width: rect.Width, Height: 1}, ' ', style)
		buffer.WriteString(rect.X+1, rect.Y+row, text, style)

	}

}
