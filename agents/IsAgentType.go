package agents

func IsAgentType(raw string) bool {

	found := false

	for _, agent_type := range AgentTypes {

		if agent_type.String() == raw {
			found = true
			break
		}

	}

	return found

}
