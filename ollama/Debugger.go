package ollama

import "exocomp/schemas"
import "bufio"
import "fmt"
import "os"
import "strings"
import "sync"

type Debugger struct {
	Prompt    string
	Session   *Session
	mutex     *sync.Mutex
	rendered  int
	resetline string
}

func NewDebugger(session *Session) *Debugger {

	resetline := ""

	for r := 0; r < len(session.Config.Model) + 10; r++ {
		resetline += " "
	}


	return &Debugger{
		Prompt:    "",
		Session:   session,
		mutex:     &sync.Mutex{},
		rendered:  0,
		resetline: resetline,
	}

}

func (debugger *Debugger) ClearLine() {

	fmt.Fprintf(os.Stdout, "\033[A")
	fmt.Fprintf(os.Stdout, "\033[K")
	fmt.Fprintf(os.Stdout, debugger.resetline)
	fmt.Fprintf(os.Stdout, "\033[B")

}

func (debugger *Debugger) Destroy() {

}

func (debugger *Debugger) InputLoop() {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		prompt := strings.TrimSpace(scanner.Text())

		if prompt != "" {

			go debugger.Session.Query(schemas.Message{
				Role:    "user",
				Content: prompt,
			})

		}

	}

}

func (debugger *Debugger) RenderLoop() {

	for {

		if debugger.rendered < len(debugger.Session.Messages) {

			next_message := debugger.Session.Messages[debugger.rendered]

			if next_message.Role == "user" {
				debugger.ClearLine()
			}

			debugger.RenderMessages(debugger.Session.Messages[debugger.rendered:])
			debugger.RenderPrompt()
			debugger.rendered = len(debugger.Session.Messages)

		}

	}

}

func (debugger *Debugger) RenderMessages(messages []schemas.Message) {

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
			color = ColorYellow
		default:
			color = ColorReset
		}

		if color != ColorReset && len(content) > 0 {

			if len(content) > 1 {

				resetline := debugger.resetline[0:len(debugger.resetline) - len(role) - 3]

				fmt.Fprintf(os.Stdout, "\r%s[%s]%s:%s\n", color, role, ColorReset, resetline)

				for _, line := range content {
					fmt.Fprintf(os.Stdout, "\r%s|%s %s\n", color, ColorReset, line)
				}

				os.Stdout.Sync()

			} else {

				resetline := ""

				if len(content[0]) < len(debugger.resetline) - 4 {
					resetline = debugger.resetline[0:len(debugger.resetline) - len(content[0]) - 4]
				}

				fmt.Fprintf(os.Stdout, "\r%s[%s]%s: %s%s\n", color, role, ColorReset, content[0], resetline)
				os.Stdout.Sync()

			}

		}

	}

}

func (debugger *Debugger) RenderPrompt() {

	model := "unknown"

	if debugger.Session != nil && debugger.Session.Config != nil {
		model = debugger.Session.Config.Model
	}

	fmt.Fprintf(os.Stdout, "\r%s[to %s]%s > ", ColorGreen, model, ColorReset)
	os.Stdout.Sync()

}
