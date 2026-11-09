package agents

func IsAgentRole(search string) bool {

	found := false

	for raw, _ := range Roles {

		if raw == search {

			found = true
			break

		}

	}

	return found

}
