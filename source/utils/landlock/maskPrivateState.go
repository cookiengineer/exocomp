//go:build linux

package landlock

import "fmt"
import "os"
import "path/filepath"

import "golang.org/x/sys/unix"

// maskPrivateState hides the sandbox's private .exocomp state folder behind an
// empty read-only tmpfs. The helper runs in a private mount namespace, so the
// mount only affects the sandboxed process and its children.
func maskPrivateState(sandbox string) error {

	private := filepath.Join(sandbox, ".exocomp")

	err0 := os.MkdirAll(private, 0700)

	if err0 == nil {

		err1 := unix.Mount("", "/", "", uintptr(unix.MS_REC|unix.MS_PRIVATE), "")

		if err1 == nil {

			err2 := unix.Mount(
				"tmpfs",
				private,
				"tmpfs",
				uintptr(unix.MS_RDONLY|unix.MS_NOSUID|unix.MS_NODEV),
				"size=1m",
			)

			if err2 == nil {
				return DropPrivileges()
			} else {
				return fmt.Errorf("landlock.MaskPrivateState: Cannot mask private state folder.")
			}

		} else {
			return fmt.Errorf("landlock.MaskPrivateState: Cannot make mount namespace private.")
		}

	} else {
		return err0
	}

}
