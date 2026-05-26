package agents

import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"

func SetRole(role string, agent *types.Agent) bool {

	role = utils_fmt.FormatAgentRole(role)

	if role != "" && agent != nil {

		if role == agent.Role {

			Roles[role] = agent

			return true

		}

	}

	return false

}
