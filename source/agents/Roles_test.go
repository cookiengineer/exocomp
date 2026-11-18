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
		"reverser",
		"summarizer",
		"tester",
		"threathunter",
		"webresearcher",
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
		"architect", "coder", "planner", "reverser", "summarizer", "tester", "webresearcher",
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
		"findings.Export",
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

	for _, role := range []string{"webscanner", "threathunter"} {

		agent, ok := Roles[role]

		if ok == true && slices.Contains(agent.AllowedTools, "findings.Import") == false {
			t.Errorf("Expected role %q to allow findings.Import", role)
		}

	}

	reproducer_methods := []string{
		"findings.BindEvidence",
		"findings.Export",
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

func TestRoles_AssetsAccess(t *testing.T) {

	importers := []string{"webscanner", "threathunter"}
	importer_methods := []string{
		"assets.Add",
		"assets.Export",
		"assets.Get",
		"assets.Import",
		"assets.List",
		"assets.Search",
	}

	for _, role := range importers {

		agent, ok := Roles[role]

		if ok == true {

			for _, method := range importer_methods {

				if slices.Contains(agent.AllowedTools, method) == false {
					t.Errorf("Expected role %q to allow %q", role, method)
				}

			}

		} else {
			t.Errorf("Expected role %q to be present", role)
		}

	}

	planner_methods := []string{"assets.Coverage", "assets.Get", "assets.List", "assets.Neighbors", "assets.Search"}
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

func TestRoles_ReproducerContainers(t *testing.T) {

	reproducer, ok := Roles["reproducer"]

	if ok == false {
		t.Fatalf("Expected role %q to be present", "reproducer")
	}

	container_methods := []string{
		"containers.Create",
		"containers.Exec",
		"containers.Images",
		"containers.List",
		"containers.Logs",
		"containers.Remove",
		"containers.Start",
		"containers.Stat",
		"containers.Stop",
	}

	for _, method := range container_methods {

		if slices.Contains(reproducer.AllowedTools, method) == false {
			t.Errorf("Expected role %q to allow %q", "reproducer", method)
		}

	}

	if slices.Contains(reproducer.AllowedImages, "docker.io/library/golang:1.27") == false {
		t.Errorf("Expected reproducer to allow the golang test image")
	}

	if slices.Contains(reproducer.AllowedPrograms, "podman") == true {
		t.Errorf("Expected reproducer to not allow podman as a program")
	}

}

func TestRoles_PromptsDoNotMentionInternalState(t *testing.T) {

	for role, agent := range Roles {

		if strings.Contains(agent.Prompt, ".exocomp") == true {
			t.Errorf("Expected role %q prompt not to mention internal state", role)
		}

		for _, tool := range agent.AllowedTools {

			if strings.Contains(tool, ".exocomp") == true {
				t.Errorf("Expected role %q allowed tool %q not to mention internal state", role, tool)
			}

		}

	}

}

func TestRoles_NonExecutionRolesHaveNoPrograms(t *testing.T) {

	for _, role := range []string{"planner", "architect", "summarizer", "webresearcher"} {

		agent, ok := Roles[role]

		if ok == false {
			t.Fatalf("Expected role %q to be present", role)
		}

		if len(agent.AllowedPrograms) != 0 {
			t.Errorf("Expected role %q to have no allowed programs, got %v", role, agent.AllowedPrograms)
		}

		for _, tool := range agent.AllowedTools {

			if strings.HasPrefix(tool, "programs.") == true {
				t.Errorf("Expected role %q to not allow %q", role, tool)
			}

		}

	}

}

func TestRoles_WebresearcherIsReadOnly(t *testing.T) {

	webresearcher, ok := Roles["webresearcher"]

	if ok == false {
		t.Fatalf("Expected role %q to be present", "webresearcher")
	}

	if len(webresearcher.AllowedPrograms) != 0 {
		t.Errorf("Expected role %q to have no allowed programs, got %v", "webresearcher", webresearcher.AllowedPrograms)
	}

	required_tools := []string{"files.Read", "files.Search", "websites.Fetch", "websites.List", "websites.Stat", "agents.Quit"}

	for _, tool := range required_tools {

		if slices.Contains(webresearcher.AllowedTools, tool) == false {
			t.Errorf("Expected role %q to allow tool %q", "webresearcher", tool)
		}

	}

	for _, tool := range webresearcher.AllowedTools {

		if strings.HasPrefix(tool, "programs.") == true {
			t.Errorf("Expected role %q to not allow %q", "webresearcher", tool)
		}

		if strings.HasPrefix(tool, "shell.") == true {
			t.Errorf("Expected role %q to not allow %q", "webresearcher", tool)
		}

		if strings.HasPrefix(tool, "assets.") == true {
			t.Errorf("Expected role %q to not allow %q", "webresearcher", tool)
		}

		if strings.HasPrefix(tool, "findings.") == true {
			t.Errorf("Expected role %q to not allow %q", "webresearcher", tool)
		}

	}

	if slices.Contains(webresearcher.AllowedTools, "files.Write") == true {
		t.Errorf("Expected role %q to not allow files.Write", "webresearcher")
	}

}

func TestRoles_ExploiterAndReverserReadAssets(t *testing.T) {

	for _, role := range []string{"exploiter", "reverser"} {

		agent, ok := Roles[role]

		if ok == false {
			t.Fatalf("Expected role %q to be present", role)
		}

		for _, method := range []string{"assets.Export", "assets.Get", "assets.List", "assets.Search"} {

			if slices.Contains(agent.AllowedTools, method) == false {
				t.Errorf("Expected role %q to allow %q", role, method)
			}

		}

	}

}
