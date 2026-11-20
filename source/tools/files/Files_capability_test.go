package files

import "exocomp/types"
import utils_capabilities "exocomp/utils/capabilities"
import "os"
import "path/filepath"
import "strings"
import "testing"

func TestFiles_CapabilityReadWrite(t *testing.T) {

	outside_folder := t.TempDir()
	outside_file := filepath.Join(outside_folder, "outside.txt")

	err0 := os.WriteFile(outside_file, []byte("outside content"), 0644)

	if err0 != nil {
		t.Fatalf("Cannot create outside file: %v", err0)
	}

	playground := t.TempDir()
	sandbox := filepath.Join(playground, "sandbox")
	capabilities := utils_capabilities.NewSet()

	tool := NewFilesWithCapabilities([]string{"Read", "Write"}, playground, sandbox, capabilities)

	_, err1 := tool.Read(outside_file)

	if err1 == nil {
		t.Errorf("Expected read without capability to fail")
	}

	err2 := capabilities.Grant([]types.Capability{
		{Type: types.CapabilityFilesRead, Scope: outside_folder},
	})

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	result, err3 := tool.Read(outside_file)

	if err3 == nil {

		if strings.Contains(result, "outside content") == false {
			t.Errorf("Expected outside content, got %s", result)
		}

	} else {
		t.Errorf("Expected %v to be nil", err3)
	}

	_, err4 := tool.Write(outside_file, "changed content")

	if err4 == nil {
		t.Errorf("Expected write without capability to fail")
	}

	err5 := capabilities.Grant([]types.Capability{
		{Type: types.CapabilityFilesWrite, Scope: outside_folder},
	})

	if err5 != nil {
		t.Fatalf("Expected %v to be nil", err5)
	}

	_, err6 := tool.Write(outside_file, "changed content")

	if err6 == nil {

		content, err7 := os.ReadFile(outside_file)

		if err7 == nil {

			if strings.Contains(string(content), "changed content") == false {
				t.Errorf("Expected changed content, got %s", string(content))
			}

		} else {
			t.Errorf("Expected %v to be nil", err7)
		}

	} else {
		t.Errorf("Expected %v to be nil", err6)
	}

}
