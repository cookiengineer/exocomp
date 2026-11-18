package agents

import "slices"
import "strings"
import "testing"

func TestRoles_ContainsExpectedRoles(t *testing.T) {

	expected := []string{
		"architect",
		"coder",
		"exploiter",
		"planner",
		"reproducer",
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

func TestRoles_AllowedProgramsAreBareNames(t *testing.T) {

	for role, agent := range Roles {

		for _, program := range agent.AllowedPrograms {

			if program == "" {
				t.Errorf("Expected role %q to not allow an empty program name", role)
			} else if strings.ContainsAny(program, "/\\") == true {
				t.Errorf("Expected role %q allowed program %q to not contain a path separator", role, program)
			} else if strings.ContainsAny(program, " \t\r\n") == true {
				t.Errorf("Expected role %q allowed program %q to not contain whitespace", role, program)
			}

		}

	}

}

func TestRoles_AllowedProgramsHaveExecutionTools(t *testing.T) {

	required := []string{"programs.Execute", "programs.List", "programs.Stat"}

	for role, agent := range Roles {

		if len(agent.AllowedPrograms) > 0 {

			for _, tool := range required {

				if slices.Contains(agent.AllowedTools, tool) == false {
					t.Errorf("Expected role %q to allow %q for its allowed programs", role, tool)
				}

			}

		}

	}

}

func TestRoles_ContainsScannerPrograms(t *testing.T) {

	expected := []struct {
		role     string
		programs []string
	}{
		{
			role: "webscanner",
			programs: []string{
				"alterx", "amass", "asnmap", "cdncheck", "dnsx", "gobuster",
				"httpx", "interactsh-client", "katana", "mapcidr", "naabu",
				"notify", "nuclei", "proxify", "subfinder", "tlsx", "urlfinder", "vulnx",
			},
		},
		{
			role: "threathunter",
			programs: []string{
				"asnmap", "cdncheck", "cloudlist", "dnsx", "interactsh-client",
				"mapcidr", "naabu", "notify", "tlsx", "vulnx",
			},
		},
		{
			role: "researcher",
			programs: []string{
				"alterx", "amass", "chaos", "dnsx", "subfinder", "uncover",
				"urlfinder", "vulnx",
			},
		},
		{
			role: "exploiter",
			programs: []string{
				"interactsh-client", "proxify", "vulnx",
			},
		},
	}

	for _, entry := range expected {

		agent, ok := Roles[entry.role]

		if ok == true {

			for _, program := range entry.programs {

				if slices.Contains(agent.AllowedPrograms, program) == false {
					t.Errorf("Expected role %q to allow scanner %q", entry.role, program)
				}

			}

		} else {
			t.Errorf("Expected role %q to be present", entry.role)
		}

	}

}

func TestRoles_NonScannerRolesDoNotAllowScanners(t *testing.T) {

	scanners := []string{
		"alterx", "asnmap", "cdncheck", "chaos", "cloudlist", "dnsx",
		"httpx", "interactsh-client", "katana", "mapcidr", "naabu", "notify",
		"nuclei", "proxify", "subfinder", "tlsx", "uncover", "urlfinder", "vulnx",
	}

	roles := []string{
		"architect", "coder", "planner", "reverser", "summarizer", "tester",
	}

	for _, role := range roles {

		agent, ok := Roles[role]

		if ok == true {

			for _, scanner := range scanners {

				if slices.Contains(agent.AllowedPrograms, scanner) == true {
					t.Errorf("Expected role %q to not allow scanner %q", role, scanner)
				}

			}

		} else {
			t.Errorf("Expected role %q to be present", role)
		}

	}

}

func TestRoles_FindingsAccess(t *testing.T) {

	writers := []string{"webscanner", "threathunter", "exploiter"}
	writer_methods := []string{
		"findings.Add",
		"findings.BindEvidence",
		"findings.Get",
		"findings.List",
		"findings.Report",
		"findings.Retest",
		"findings.Search",
		"findings.Update",
		"findings.Verify",
	}

	for _, role := range writers {

		agent, ok := Roles[role]

		if ok == true {

			for _, method := range writer_methods {

				if slices.Contains(agent.AllowedTools, method) == false {
					t.Errorf("Expected role %q to allow %q", role, method)
				}

			}

		} else {
			t.Errorf("Expected role %q to be present", role)
		}

	}

	reproducer_methods := []string{
		"findings.BindEvidence",
		"findings.Get",
		"findings.List",
		"findings.Report",
		"findings.Retest",
		"findings.Search",
		"findings.Verify",
	}

	reproducer, ok := Roles["reproducer"]

	if ok == true {

		for _, method := range reproducer_methods {

			if slices.Contains(reproducer.AllowedTools, method) == false {
				t.Errorf("Expected role %q to allow %q", "reproducer", method)
			}

		}

		if slices.Contains(reproducer.AllowedTools, "findings.Add") == true {
			t.Errorf("Expected role %q to not allow findings.Add", "reproducer")
		}

		if slices.Contains(reproducer.AllowedTools, "findings.Update") == true {
			t.Errorf("Expected role %q to not allow findings.Update", "reproducer")
		}

	} else {
		t.Errorf("Expected role %q to be present", "reproducer")
	}

	planner_methods := []string{"findings.Get", "findings.List", "findings.Report", "findings.Search"}

	planner, ok := Roles["planner"]

	if ok == true {

		for _, method := range planner_methods {

			if slices.Contains(planner.AllowedTools, method) == false {
				t.Errorf("Expected role %q to allow %q", "planner", method)
			}

		}

	} else {
		t.Errorf("Expected role %q to be present", "planner")
	}

}
