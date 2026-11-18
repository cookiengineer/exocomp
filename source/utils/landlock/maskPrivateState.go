//go:build linux

package landlock

import "fmt"
import "os"
import "path/filepath"

import "golang.org/x/sys/unix"

const (
	securebitsNoRoot       = 1 << 0
	securebitsNoRootLocked = 1 << 1
)

// maskPrivateState hides the sandbox's private .exocomp state folder behind an
// empty read-only tmpfs. The helper runs in a private mount namespace, so the
// mount only affects the sandboxed process and its children.
func maskPrivateState(sandbox string) error {

	private := filepath.Join(sandbox, ".exocomp")

	err0 := os.MkdirAll(private, 0700)

	if err0 != nil {
		return err0
	}

	err1 := unix.Mount("", "/", "", uintptr(unix.MS_REC|unix.MS_PRIVATE), "")

	if err1 != nil {
		return fmt.Errorf("Cannot make mount namespace private.")
	}

	err2 := unix.Mount(
		"tmpfs",
		private,
		"tmpfs",
		uintptr(unix.MS_RDONLY|unix.MS_NOSUID|unix.MS_NODEV),
		"size=1m",
	)

	if err2 != nil {
		return fmt.Errorf("Cannot mask private state folder.")
	}

	return dropPrivileges()

}

// dropPrivileges clears the capabilities gained from entering the user
// namespace. The target must not run with CAP_DAC_OVERRIDE or similar, because
// that would bypass ordinary filesystem permissions inside the sandbox.
func dropPrivileges() error {

	err0 := unix.Prctl(unix.PR_SET_SECUREBITS, uintptr(securebitsNoRoot|securebitsNoRootLocked), 0, 0, 0)

	if err0 != nil {
		return fmt.Errorf("Cannot set secure bits.")
	}

	for capability := 0; capability <= int(unix.CAP_LAST_CAP); capability++ {
		_ = unix.Prctl(unix.PR_CAPBSET_DROP, uintptr(capability), 0, 0, 0)
	}

	data := [2]unix.CapUserData{}
	err1 := unix.Capset(&unix.CapUserHeader{Version: unix.LINUX_CAPABILITY_VERSION_3}, &data[0])

	if err1 != nil {
		return fmt.Errorf("Cannot drop capabilities.")
	}

	return nil

}
