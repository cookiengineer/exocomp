package tty

import "fmt"
import "strings"
import "sync"

type Screen struct {
	terminal          *Terminal
	mutex             *sync.Mutex
	buffer            *Buffer
	previous          *Buffer
	cursorX           int
	cursorY           int
	cursorVisible     bool
	profile           ColorProfile
	width             int
	height            int
	lastCursorX       int
	lastCursorY       int
	lastCursorVisible bool
	hasCursor         bool
}

func NewScreen(terminal *Terminal, width int, height int) *Screen {

	return &Screen{
		terminal:      terminal,
		mutex:         &sync.Mutex{},
		buffer:        NewBuffer(width, height),
		previous:      NewBuffer(width, height),
		cursorX:       0,
		cursorY:       0,
		cursorVisible: true,
		profile:       terminal.Profile,
		width:         width,
		height:        height,
	}

}

func (screen *Screen) Buffer() *Buffer {
	return screen.buffer
}

func (screen *Screen) Width() int {
	return screen.buffer.Width
}

func (screen *Screen) Height() int {
	return screen.buffer.Height
}

func (screen *Screen) Resize(width int, height int) {

	screen.mutex.Lock()

	screen.buffer.Resize(width, height)
	screen.previous.Resize(width, height)
	screen.previous.Clear(Style{})
	screen.width = width
	screen.height = height

	screen.mutex.Unlock()

}

func (screen *Screen) SetCursor(x int, y int, visible bool) {

	screen.cursorX = x
	screen.cursorY = y
	screen.cursorVisible = visible

}

func (screen *Screen) Flush() error {

	screen.mutex.Lock()
	defer screen.mutex.Unlock()

	if screen.buffer.Width != screen.previous.Width || screen.buffer.Height != screen.previous.Height {
		screen.previous = NewBuffer(screen.buffer.Width, screen.buffer.Height)
	}

	changed := false
	output := strings.Builder{}
	output.WriteString("\033[?2026h")

	current_style := Style{}
	has_style := false
	current_x := -1
	current_y := -1

	for y := 0; y < screen.buffer.Height; y++ {

		for x := 0; x < screen.buffer.Width; {

			cell := screen.buffer.Get(x, y)
			previous := screen.previous.Get(x, y)

			if cell.Width == 0 {

				x++
				continue

			}

			if cell.Equals(previous) {

				x += max(1, cell.Width)
				continue

			}

			changed = true

			if current_x != x || current_y != y {

				output.WriteString(fmt.Sprintf("\033[%d;%dH", y+1, x+1))
				current_x = x
				current_y = y

			}

			if has_style == false || current_style.Equals(cell.Style) == false {

				output.WriteString(cell.Style.Sequence(screen.profile))
				current_style = cell.Style
				has_style = true

			}

			if cell.Text == "" {
				output.WriteString(" ")
			} else {
				output.WriteString(cell.Text)
			}

			x += max(1, cell.Width)
			current_x = x

		}

	}

	if changed == false && screen.hasCursor == true && screen.lastCursorX == screen.cursorX && screen.lastCursorY == screen.cursorY && screen.lastCursorVisible == screen.cursorVisible {
		return nil
	}

	if screen.profile != ProfileMono {
		output.WriteString("\033[0m")
	}

	if screen.cursorVisible == true {

		output.WriteString(fmt.Sprintf("\033[%d;%dH", screen.cursorY+1, screen.cursorX+1))
		output.WriteString("\033[?25h")

	} else {
		output.WriteString("\033[?25l")
	}

	output.WriteString("\033[?2026l")

	if output.Len() > 0 {
		screen.terminal.Write(output.String())
	}

	screen.previous = screen.buffer.Clone()
	screen.lastCursorX = screen.cursorX
	screen.lastCursorY = screen.cursorY
	screen.lastCursorVisible = screen.cursorVisible
	screen.hasCursor = true

	return nil

}
