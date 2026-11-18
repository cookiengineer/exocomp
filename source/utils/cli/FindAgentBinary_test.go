package cli

import "os"
import "path/filepath"
import "runtime"
import "strings"
import "testing"

func writeFakeAgent(t *testing.T, directory string) string {

	t.Helper()

	name := "exocomp-agent"

	if runtime.GOOS == "windows" {
		name = name + ".exe"
	}

	path := filepath.Join(directory, name)
	err := os.WriteFile(path, []byte("fake"), 0755)

	if err == nil {
		return path
	} else {
		t.Fatalf("Expected fake agent to be written, got %v", err)
		return ""
	}

}

func TestFindAgentBinary_EnvOverride(t *testing.T) {

	directory := t.TempDir()
	override := writeFakeAgent(t, directory)
	other := t.TempDir()

	t.Setenv("EXOCOMP_AGENT", override)
	t.Setenv("PATH", "")

	result, err0 := findAgentBinaryIn(other)

	if err0 == nil {

		expected, _ := filepath.Abs(override)

		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

	} else {
		t.Errorf("Expected %v to be nil", err0)
	}

}

func TestFindAgentBinary_Sibling(t *testing.T) {

	directory := t.TempDir()
	sibling := writeFakeAgent(t, directory)

	t.Setenv("EXOCOMP_AGENT", "")
	t.Setenv("PATH", "")

	result, err0 := findAgentBinaryIn(directory)

	if err0 == nil {

		expected, _ := filepath.Abs(sibling)

		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

	} else {
		t.Errorf("Expected %v to be nil", err0)
	}

}

func TestFindAgentBinary_Path(t *testing.T) {

	directory := t.TempDir()
	search := t.TempDir()
	path := writeFakeAgent(t, search)

	t.Setenv("EXOCOMP_AGENT", "")
	t.Setenv("PATH", search)

	result, err0 := findAgentBinaryIn(directory)

	if err0 == nil {

		expected, _ := filepath.Abs(path)

		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

	} else {
		t.Errorf("Expected %v to be nil", err0)
	}

}

func TestFindAgentBinary_Missing(t *testing.T) {

	directory := t.TempDir()

	t.Setenv("EXOCOMP_AGENT", "")
	t.Setenv("PATH", "")

	result, err0 := findAgentBinaryIn(directory)

	if err0 == nil {
		t.Errorf("Expected %v to be not nil", err0)
	} else if strings.Contains(err0.Error(), "not installed") == false {
		t.Errorf("Expected not installed error, got %v", err0)
	} else if strings.Contains(err0.Error(), directory) == true {
		t.Errorf("Expected error to omit the host path, got %v", err0)
	} else if result != "" {
		t.Errorf("Expected empty result, got %q", result)
	}

}
