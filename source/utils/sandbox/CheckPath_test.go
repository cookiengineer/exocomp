package sandbox

import "os"
import "path/filepath"
import "testing"

func TestCheckPath(t *testing.T) {

	sandbox := t.TempDir()
	subfolder := filepath.Join(sandbox, "sub")
	err0 := os.MkdirAll(subfolder, 0755)

	if err0 != nil {
		t.Fatalf("Cannot create test folder: %v", err0)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"", "."},
		{".", "."},
		{"./", "."},
		{sandbox, "."},
		{"file.txt", "file.txt"},
		{"./file.txt", "file.txt"},
		{filepath.Join(sandbox, "sub", "file.txt"), filepath.Join("sub", "file.txt")},
	}

	for _, test := range tests {

		result, err1 := CheckPath(sandbox, test.input)

		if err1 != nil {
			t.Errorf("Expected %q to be allowed, got %v", test.input, err1)
		} else if result != test.expected {
			t.Errorf("Expected %q, got %q", test.expected, result)
		}

	}

	_, err2 := CheckPath(sandbox, "../outside.txt")

	if err2 == nil {
		t.Errorf("Expected parent path to be rejected")
	}

}

func TestCheckSymlinks(t *testing.T) {

	sandbox := t.TempDir()
	file := filepath.Join(sandbox, "file.txt")

	err0 := os.WriteFile(file, []byte("test"), 0644)

	if err0 != nil {
		t.Fatalf("Cannot create test file: %v", err0)
	}

	createTestSymlink(t, file, filepath.Join(sandbox, "link.txt"))

	err1 := CheckSymlinks(sandbox, "link.txt")

	if err1 == nil {
		t.Errorf("Expected symlink to be rejected")
	}

	err2 := CheckSymlinks(sandbox, "file.txt")

	if err2 != nil {
		t.Errorf("Expected regular file to be allowed, got %v", err2)
	}

}
