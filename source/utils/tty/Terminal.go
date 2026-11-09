package tty

import "encoding/base64"
import "fmt"
import "os"
import "sync"
import "golang.org/x/term"

type Terminal struct {
	Input   *os.File
	Output  *os.File
	Profile ColorProfile
	mutex   *sync.Mutex
	state   *term.State
	width   int
	height  int
	entered bool
}

func NewTerminal() *Terminal {

	return &Terminal{
		Input:   os.Stdin,
		Output:  os.Stdout,
		Profile: DetectColorProfile(),
		mutex:   &sync.Mutex{},
		width:   -1,
		height:  -1,
	}

}

func (terminal *Terminal) Enter() error {

	terminal.mutex.Lock()
	defer terminal.mutex.Unlock()

	if terminal.entered == true {
		return nil
	}

	state, err := term.MakeRaw(int(terminal.Input.Fd()))

	if err == nil {

		terminal.state = state
		terminal.entered = true
		terminal.write("\033[?1049h")
		terminal.write("\033[?25l")
		terminal.write("\033[?2004h")
		terminal.write("\033[2J\033[H")
		terminal.Output.Sync()

		return nil

	} else {
		return fmt.Errorf("tty.Terminal.Enter: %s", err.Error())
	}

}

func (terminal *Terminal) Leave() {

	terminal.mutex.Lock()
	defer terminal.mutex.Unlock()

	if terminal.entered == false {
		return
	}

	terminal.write("\033[?2004l")
	terminal.write("\033[?25h")
	terminal.write("\033[?1049l")
	terminal.Output.Sync()

	if terminal.state != nil {

		term.Restore(int(terminal.Input.Fd()), terminal.state)
		terminal.state = nil

	}

	terminal.entered = false

}

func (terminal *Terminal) Size() (int, int) {

	width, height, err := term.GetSize(int(terminal.Output.Fd()))

	if err == nil && width > 0 && height > 0 {

		terminal.width = width
		terminal.height = height

		return width, height

	} else {

		if terminal.width < 0 || terminal.height < 0 {

			terminal.width = 80
			terminal.height = 24

		}

		return terminal.width, terminal.height

	}

}

func (terminal *Terminal) Write(text string) {

	terminal.mutex.Lock()
	terminal.write(text)
	terminal.Output.Sync()
	terminal.mutex.Unlock()

}

func (terminal *Terminal) Copy(text string) error {

	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	terminal.Write("\033]52;c;" + encoded + "\a")

	return nil

}

func (terminal *Terminal) write(text string) {

	terminal.Output.WriteString(text)

}
