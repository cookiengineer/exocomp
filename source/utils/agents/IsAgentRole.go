package agents

import "exocomp/agents"

func IsAgentRole(search string) bool {

	found := false

	for raw, _ := range agents.Roles {

		if raw == search {
			found = true
			break
		}

	}

	return found

}
