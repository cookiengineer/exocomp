package testutil

import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "strings"
import "sync"
import "testing"

var (
	fakePodmanOnce sync.Once
	fakePodmanPath string
	fakePodmanErr  error
)

// BuildFakePodman compiles the fake podman binary once and returns its path.
func BuildFakePodman(t *testing.T) string {

	t.Helper()

	fakePodmanOnce.Do(func() {

		root := sourceRoot()

		if root == "" {

			fakePodmanErr = fmt.Errorf("Cannot locate the source folder")
			return

		}

		output := filepath.Join(os.TempDir(), "exocomp-fakepodman-test")
		cmd := exec.Command("go", "build", "-o", output, "./tools/testdata/fakepodman")
		cmd.Dir = root

		data, err := cmd.CombinedOutput()

		if err != nil {

			fakePodmanErr = fmt.Errorf("Cannot build fakepodman: %v: %s", err, strings.TrimSpace(string(data)))
			return

		}

		fakePodmanPath = output

	})

	if fakePodmanErr != nil {
		t.Fatalf("%v", fakePodmanErr)
	}

	return fakePodmanPath

}
