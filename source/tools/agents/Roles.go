package agents

import agents_roles "exocomp/agents"
import "fmt"
import "sort"
import "strings"

func (tool *Agents) Roles() (string, error) {

	lines := make([]string, 0)

	for _, template := range agents_roles.Roles {

		if template.Role != "planner" {
			lines = append(lines, fmt.Sprintf("- Role: \"%s\", Description: %s", template.Role, template.Description))
		}

	}

	sort.Strings(lines)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("agents.Roles: %d available agent roles.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}
