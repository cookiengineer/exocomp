package testutil

import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "strings"
import "sync"
import "testing"

var (
	sandboxExecOnce sync.Once
	sandboxExecPath string
	sandboxExecErr  error
)

// BuildSandboxExec compiles the exocomp-sandbox helper once and returns its
// path.
func BuildSandboxExec(t *testing.T) string {

	t.Helper()

	sandboxExecOnce.Do(func() {

		root := sourceRoot()

		if root == "" {
			sandboxExecErr = fmt.Errorf("Cannot locate the source folder")
			return
		}

		output := filepath.Join(os.TempDir(), "exocomp-sandbox-test")
		cmd := exec.Command("go", "build", "-o", output, "./cmds/exocomp-sandbox")
		cmd.Dir = root

		data, err := cmd.CombinedOutput()

		if err != nil {
			sandboxExecErr = fmt.Errorf("Cannot build exocomp-sandbox: %v: %s", err, strings.TrimSpace(string(data)))
			return
		}

		sandboxExecPath = output

	})

	if sandboxExecErr != nil {
		t.Fatalf("%v", sandboxExecErr)
	}

	return sandboxExecPath

}
