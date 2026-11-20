package testutil

import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "strings"
import "sync"
import "testing"

var (
	fakeAgentOnce sync.Once
	fakeAgentPath string
	fakeAgentErr  error
)

// BuildFakeAgent compiles the fake agent binary once and returns its path.
func BuildFakeAgent(t *testing.T) string {

	t.Helper()

	fakeAgentOnce.Do(func() {

		root := sourceRoot()

		if root == "" {
			fakeAgentErr = fmt.Errorf("Cannot locate the source folder")
			return
		}

		output := filepath.Join(os.TempDir(), "exocomp-fake-agent")
		cmd := exec.Command("go", "build", "-o", output, "./tools/testdata/fakeagent")
		cmd.Dir = root

		data, err := cmd.CombinedOutput()

		if err != nil {
			fakeAgentErr = fmt.Errorf("Cannot build fake agent: %v: %s", err, strings.TrimSpace(string(data)))
			return
		}

		fakeAgentPath = output

	})

	if fakeAgentErr != nil {
		t.Fatalf("%v", fakeAgentErr)
	}

	return fakeAgentPath

}
