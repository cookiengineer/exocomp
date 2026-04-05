package tty

import "exocomp/ollama"
import "exocomp/schemas"
import "bufio"
import "fmt"
import "os"
import "sort"
import "strings"
import "sync"

type Renderer struct {
	Prompt    string
	Session   *ollama.Session
	mutex     *sync.RWMutex
	rendered  int
	resetline string
}

func NewRenderer(session *ollama.Session) *Renderer {

	resetline := ""

	for r := 0; r < len(session.Config.Model) + 10; r++ {
		resetline += " "
	}


	return &Renderer{
		Prompt:    "",
		Session:   session,
		mutex:     &sync.RWMutex{},
		rendered:  0,
		resetline: resetline,
	}

}

func (renderer *Renderer) ClearLine() {

	fmt.Fprintf(os.Stdout, "\033[A")
	fmt.Fprintf(os.Stdout, "\033[K")
	fmt.Fprintf(os.Stdout, renderer.resetline)
	fmt.Fprintf(os.Stdout, "\033[B")

}

func (renderer *Renderer) Destroy() {

}

func (renderer *Renderer) InputLoop() {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		prompt := strings.TrimSpace(scanner.Text())

		if prompt != "" {

			go renderer.Session.Query(schemas.Message{
				Role:    "user",
				Content: prompt,
			})

		}

	}

}

func (renderer *Renderer) RenderLoop() {

	tools := make([]string, 0)

	for _, tool := range renderer.Session.Agent.Tools {
		tools = append(tools, tool)
	}

	sort.Strings(tools)

	info_agent := fmt.Sprintf("Agent: %s | %s | %.2f", renderer.Session.Agent.Type, renderer.Session.Agent.Model, renderer.Session.Agent.Temperature)
	info_tools := fmt.Sprintf("Tools: %s", strings.Join(tools, ", "))

	fmt.Fprintf(os.Stdout, "\r%s[exocomp]%s:\n", ColorYellow, ColorReset)
	fmt.Fprintf(os.Stdout, "\r%s|%s %s\n", ColorYellow, ColorReset, info_agent)
	fmt.Fprintf(os.Stdout, "\r%s|%s %s\n", ColorYellow, ColorReset, info_tools)
	fmt.Fprintf(os.Stdout, "\n")
	os.Stdout.Sync()

	for {

		renderer.mutex.RLock()
		from := renderer.rendered
		renderer.mutex.RUnlock()

		messages := renderer.Session.GetMessages(from)

		if len(messages) > 0 {

			if messages[0] != nil && messages[0].Role == "user" {
				renderer.ClearLine()
			}

			renderer.RenderMessages(messages)
			renderer.RenderPrompt()

			renderer.mutex.Lock()
			renderer.rendered += len(messages)
			renderer.mutex.Unlock()

		} else {
			continue
		}

	}

}

func (renderer *Renderer) RenderMessages(messages []*schemas.Message) {

	for _, message := range messages {

		if message == nil {
			continue
		}

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

				resetline := renderer.resetline[0:len(renderer.resetline) - len(role) - 3]

				fmt.Fprintf(os.Stdout, "\r%s[%s]%s:%s\n", color, role, ColorReset, resetline)

				for _, line := range content {
					fmt.Fprintf(os.Stdout, "\r%s|%s %s\n", color, ColorReset, line)
				}

				fmt.Fprintf(os.Stdout, "\n")
				os.Stdout.Sync()

			} else if len(content) == 1 {

				resetline := ""

				if len(content[0]) < len(renderer.resetline) - 4 {
					resetline = renderer.resetline[0:len(renderer.resetline) - len(content[0]) - 4]
				}

				fmt.Fprintf(os.Stdout, "\r%s[%s]%s: %s%s\n", color, role, ColorReset, content[0], resetline)

				fmt.Fprintf(os.Stdout, "\n")
				os.Stdout.Sync()

			}

		}

	}

}

func (renderer *Renderer) RenderPrompt() {

	model := "unknown"

	if renderer.Session != nil && renderer.Session.Config != nil {
		model = renderer.Session.Config.Model
	}

	fmt.Fprintf(os.Stdout, "\r%s[to %s]%s > ", ColorGreen, model, ColorReset)
	os.Stdout.Sync()

}
