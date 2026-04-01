package agents

func IsAgentType(raw string) bool {

	switch raw {
	case "architect":
		return true
	case "coder":
		return true
	case "tester":
		return true
	case "manager":
		return true
	default:
		return false
	}

}
