package agents

import agents_roles "exocomp/agents"
import utils_fmt "exocomp/utils/fmt"
import "fmt"

func (tool *Agents) Switch(role string) (string, error) {

	role = utils_fmt.FormatAgentRole(role)

	template, ok := agents_roles.Roles[role]

	if ok == true && template != nil {

		if tool.OnSwitch != nil {
			return tool.OnSwitch(role)
		} else {
			return "", fmt.Errorf("agents.Switch: Agent profiles can only be switched in the main session.")
		}

	} else {
		return "", fmt.Errorf("agents.Switch: Unknown agent role \"%s\".", role)
	}

}
