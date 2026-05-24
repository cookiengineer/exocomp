package agents

import "exocomp/types"

func SetRole(role string, agent *types.Agent) bool {

	if role != "" && agent != nil {

		if role == agent.Role {

			Roles[role] = agent

			return true

		}

	}

	return false

}
