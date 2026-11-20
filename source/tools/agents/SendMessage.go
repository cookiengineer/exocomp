package agents

import "exocomp/schemas"
import utils_fmt "exocomp/utils/fmt"
import "encoding/json"
import "fmt"
import "io"
import "strings"

func (tool *Agents) SendMessage(name string, content string) (string, error) {

	sanitized_name := utils_fmt.FormatAgentName(name)
	message := strings.TrimSpace(content)

	if message != "" {

		tool.mutex.Lock()
		process, running := tool.processes[sanitized_name]
		tool.mutex.Unlock()

		if running == true && process != nil && process.stdin != nil {

			payload, err0 := json.Marshal(schemas.Message{
				Role:    "user",
				Content: message,
			})

			if err0 == nil {

				process.write_mutex.Lock()
				_, err1 := io.WriteString(process.stdin, string(payload)+"\n")
				process.write_mutex.Unlock()

				if err1 == nil {
					return fmt.Sprintf("agents.SendMessage: Message sent to agent \"%s\".", sanitized_name), nil
				} else {
					return "", fmt.Errorf("agents.SendMessage: Cannot write to agent \"%s\".", sanitized_name)
				}

			} else {
				return "", fmt.Errorf("agents.SendMessage: Cannot encode message for agent \"%s\".", sanitized_name)
			}

		} else {
			return "", fmt.Errorf("agents.SendMessage: Agent \"%s\" is not working.", sanitized_name)
		}

	} else {
		return "", fmt.Errorf("agents.SendMessage: Message must not be empty.")
	}

}
