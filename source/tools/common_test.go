package tools

import "os"
import "testing"

func createTestSymlink(t *testing.T, oldname string, newname string) {

	t.Helper()

	err := os.Symlink(oldname, newname)

	if err != nil {
		t.Skipf("Cannot create symlink %q -> %q: %v", newname, oldname, err)
	}

}

