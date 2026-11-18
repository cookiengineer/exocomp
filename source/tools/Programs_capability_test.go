package tools

import "exocomp/types"
import utils_capabilities "exocomp/utils/capabilities"
import "os"
import "path/filepath"
import "strings"
import "testing"

func TestPrograms_ExecuteCapabilityRead(t *testing.T) {

	requireSandbox(t)

	playground := t.TempDir()
	sandbox := filepath.Join(playground, "sandbox")
	outside_folder := t.TempDir()
	outside_file := filepath.Join(outside_folder, "outside.txt")

	err0 := os.WriteFile(outside_file, []byte("outside content"), 0644)

	if err0 != nil {
		t.Fatalf("Cannot create outside file: %v", err0)
	}

	capabilities := utils_capabilities.NewSet()
	tool := NewProgramsWithCapabilities([]string{"Execute"}, playground, sandbox, []string{"cat"}, capabilities)

	_, err1 := tool.Execute("cat", []string{outside_file})

	if err1 == nil {
		t.Errorf("Expected execution without capability to fail")
	}

	err2 := capabilities.Grant([]types.Capability{
		{Type: types.CapabilityFilesRead, Scope: outside_folder},
	})

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	result, err3 := tool.Execute("cat", []string{outside_file})

	if err3 == nil {

		if strings.Contains(result, "outside content") == false {
			t.Errorf("Expected outside content, got %s", result)
		}

	} else {
		t.Errorf("Expected %v to be nil", err3)
	}

}
