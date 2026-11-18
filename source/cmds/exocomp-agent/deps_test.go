package main

import "os/exec"
import "strings"
import "testing"

func TestDependencies_OnlyJSONL(t *testing.T) {

	cmd := exec.Command("go", "list", "-deps", ".")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Expected go list to succeed, got %v: %s", err, strings.TrimSpace(string(output)))
	}

	for _, line := range strings.Split(string(output), "\n") {

		dependency := strings.TrimSpace(line)

		if strings.HasPrefix(dependency, "exocomp/ui/") == true && dependency != "exocomp/ui/jsonl" {
			t.Errorf("Expected exocomp-agent to link exocomp/ui/jsonl only, got %q", dependency)
		}

	}

}
