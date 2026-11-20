package testutil

import "os"
import "testing"

// CreateTestSymlink creates a symlink or skips the test when the platform does
// not allow it.
func CreateTestSymlink(t *testing.T, oldname string, newname string) {

	t.Helper()

	err := os.Symlink(oldname, newname)

	if err != nil {
		t.Skipf("Cannot create symlink %q -> %q: %v", newname, oldname, err)
	}

}
