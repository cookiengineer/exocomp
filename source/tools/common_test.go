package tools

import "exocomp/utils/sandbox"
import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "runtime"
import "strings"
import "sync"
import "testing"

var (
	sandboxExecOnce sync.Once
	sandboxExecPath string
	sandboxExecErr  error
)

func buildSandboxExec(t *testing.T) string {

	sandboxExecOnce.Do(func() {

		_, filename, _, ok := runtime.Caller(0)

		if ok == false {
			sandboxExecErr = fmt.Errorf("Cannot locate test file path")
			return
		}

		tools_dir := filepath.Dir(filename)
		output := filepath.Join(os.TempDir(), "exocomp-sandbox-test")
		cmd := exec.Command("go", "build", "-o", output, "../cmds/exocomp-sandbox")
		cmd.Dir = tools_dir

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

func requireSandbox(t *testing.T) {

	t.Helper()

	if sandbox.Supported() == false {
		t.Skip("sandboxing is not available on this system")
	}

	t.Setenv("EXOCOMP_SANDBOX", buildSandboxExec(t))

}

func createTestSymlink(t *testing.T, oldname string, newname string) {

	t.Helper()

	err := os.Symlink(oldname, newname)

	if err != nil {
		t.Skipf("Cannot create symlink %q -> %q: %v", newname, oldname, err)
	}

}
