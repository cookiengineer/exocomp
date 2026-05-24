package agent

import "exocomp/agents"

func IsRole(search string) bool {

	found := false

	for role, _ := range agents.Roles {

		if role == search {
			found = true
			break
		}

	}

	return found

}
