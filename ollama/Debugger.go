package ollama

import "exocomp/schemas"
import "bufio"
import "fmt"
import "os"
import "sort"
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

	tools := make([]string, 0)

	for _, tool := range debugger.Session.Agent.Tools {
		tools = append(tools, tool)
	}

	sort.Strings(tools)

	info_agent       := fmt.Sprintf("Agent:       %s", debugger.Session.Agent.Type)
	info_model       := fmt.Sprintf("Model:       %s", debugger.Session.Config.Model)
	info_temperature := fmt.Sprintf("Temperature: %.2f", debugger.Session.Config.Temperature)
	info_tools       := fmt.Sprintf("Tools:       %s", strings.Join(tools, ", "))

	fmt.Fprintf(os.Stdout, "\r%s[exocomp]%s:\n", ColorYellow, ColorReset)
	fmt.Fprintf(os.Stdout, "\r%s|%s %s\n", ColorYellow, ColorReset, info_agent)
	fmt.Fprintf(os.Stdout, "\r%s|%s %s\n", ColorYellow, ColorReset, info_model)
	fmt.Fprintf(os.Stdout, "\r%s|%s %s\n", ColorYellow, ColorReset, info_temperature)
	fmt.Fprintf(os.Stdout, "\r%s|%s %s\n", ColorYellow, ColorReset, info_tools)
	fmt.Fprintf(os.Stdout, "\n")
	os.Stdout.Sync()

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
		limit   := len(content)

		switch message.Role {
		case "user":
			color = ColorGreen
		case "assistant":
			color = ColorBlue
		case "tool":
			color = ColorRed
			limit = 1
		case "system":
			color = ColorYellow
		default:
			color = ColorReset
		}

		if color != ColorReset && len(content) > 0 {

			if len(content) > limit {
				content = content[0:limit]
			}

			if len(content) > 1 {

				resetline := debugger.resetline[0:len(debugger.resetline) - len(role) - 3]

				fmt.Fprintf(os.Stdout, "\r%s[%s]%s:%s\n", color, role, ColorReset, resetline)

				for _, line := range content {
					fmt.Fprintf(os.Stdout, "\r%s|%s %s\n", color, ColorReset, line)
				}

				fmt.Fprintf(os.Stdout, "\n")
				os.Stdout.Sync()

			} else {

				resetline := ""

				if len(content[0]) < len(debugger.resetline) - 4 {
					resetline = debugger.resetline[0:len(debugger.resetline) - len(content[0]) - 4]
				}

				fmt.Fprintf(os.Stdout, "\r%s[%s]%s: %s%s\n", color, role, ColorReset, content[0], resetline)

				fmt.Fprintf(os.Stdout, "\n")
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
