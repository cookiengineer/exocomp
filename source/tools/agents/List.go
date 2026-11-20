package agents

import "fmt"
import "sort"
import "strings"

func (tool *Agents) List() (string, error) {

	len_content := 0

	tool.mutex.Lock()
	len_content = len(tool.contents)
	tool.mutex.Unlock()

	if len_content > 0 {

		lines := make([]string, 0)

		tool.mutex.Lock()
		for name, agent := range tool.contents {

			status := agent.Status

			if status == "" {
				status = "unknown"
			}

			lines = append(lines, fmt.Sprintf("- Name: \"%s\", Type: %s, Status: %s", name, agent.Role, status))

		}
		tool.mutex.Unlock()

		sort.Strings(lines)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("agents.List: %d agents were working for us.", len(lines)))

		for l := 0; l < len(lines); l++ {
			result = append(result, lines[l])
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("agents.List: No agents are working for us!")
	}

}
