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
}

func NewDebugger(session *Session) *Debugger {

	return &Debugger{
		Prompt:    "",
		Session:   session,
		mutex:     &sync.Mutex{},
		rendered:  0,
	}

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

			debugger.RenderMessages(debugger.Session.Messages[debugger.rendered:])
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
