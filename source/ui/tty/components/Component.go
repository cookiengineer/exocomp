package components

import utils_tty "exocomp/utils/tty"
import "time"

type Component interface {
	Width(width int) int
	Height(width int) int
	Render(buffer *utils_tty.Buffer, rect utils_tty.Rect, elapsed time.Duration)
}

var _ Component = (*Header)(nil)
var _ Component = (*Transcript)(nil)
var _ Component = (*Sidebar)(nil)
var _ Component = (*Footer)(nil)
var _ Component = (*QuestionDialog)(nil)
var _ Component = (*ChoiceDialog)(nil)
var _ Component = (*Palette)(nil)
