package agents

import "testing"

func TestRoles_ContainsExpectedRoles(t *testing.T) {

	expected := []string{
		"architect",
		"coder",
		"exploiter",
		"planner",
		"researcher",
		"reverser",
		"summarizer",
		"tester",
		"threathunter",
		"webscanner",
	}

	for _, role := range expected {

		agent, ok := Roles[role]

		if ok == false {
			t.Errorf("Expected role %q to be present", role)
		} else if agent == nil {
			t.Errorf("Expected role %q to not be nil", role)
		}

	}

}
