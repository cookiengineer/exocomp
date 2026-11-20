package testutil

import "os"
import "os/exec"
import "path/filepath"
import "strings"
import "sync"
import "testing"

var (
	fakeShellOnce sync.Once
	fakeShellPath string
	fakeShellErr  error
)

// BuildFakeShell compiles the fake shell binary once and returns its path.
func BuildFakeShell(t *testing.T) string {

	t.Helper()

	fakeShellOnce.Do(func() {

		root := sourceRoot()

		if root == "" {

			fakeShellErr = os.ErrNotExist
			return

		}

		output := filepath.Join(os.TempDir(), "exocomp-fakeshell-test")
		cmd := exec.Command("go", "build", "-o", output, "./tools/testdata/fakeshell")
		cmd.Dir = root

		data, err := cmd.CombinedOutput()

		if err == nil {
			fakeShellPath = output
		} else {

			fakeShellErr = err
			t.Logf("fake shell build output: %s", strings.TrimSpace(string(data)))

		}

	})

	if fakeShellErr != nil {
		t.Fatalf("Cannot build fake shell: %v", fakeShellErr)
	}

	return fakeShellPath

}
