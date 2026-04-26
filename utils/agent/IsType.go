package agent

import "exocomp/agents"

func IsType(raw string) bool {

	found := false

	for _, typ := range agents.Types {

		if typ == raw {
			found = true
			break
		}

	}

	return found

}
