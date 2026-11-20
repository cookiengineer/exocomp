package agents

import "fmt"
import "strings"

func (tool *Agents) Quit(message string) (string, error) {

	report := fmt.Sprintf("agents.Quit: Work Report\n%s", strings.TrimSpace(message))

	if strings.Contains(strings.ToLower(message), "my work is done") {

		if tool.OnQuit != nil {
			go tool.OnQuit(strings.TrimSpace(message), true)
		}

		return report, nil

	} else {

		if tool.OnQuit != nil {
			go tool.OnQuit(strings.TrimSpace(message), false)
		}

		return report, nil

	}

}
