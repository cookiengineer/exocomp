package shell

import "os"
import "path/filepath"
import "testing"

func TestResolveWorkdir(t *testing.T) {

	sandbox := t.TempDir()
	subfolder := filepath.Join(sandbox, "sub")
	err0 := os.MkdirAll(subfolder, 0755)

	if err0 != nil {
		t.Fatalf("Cannot create test folder: %v", err0)
	}

	result, err1 := ResolveWorkdir(sandbox, "sub")

	if err1 != nil {
		t.Fatalf("Expected %v to be nil", err1)
	}

	if result != subfolder {
		t.Errorf("Expected %q, got %q", subfolder, result)
	}

	_, err2 := ResolveWorkdir(sandbox, "missing")

	if err2 == nil {
		t.Errorf("Expected missing workdir to be rejected")
	}

	_, err3 := ResolveWorkdir(sandbox, "../outside")

	if err3 == nil {
		t.Errorf("Expected escaping workdir to be rejected")
	}

}
