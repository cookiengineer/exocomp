//go:build linux

package landlock

import "fmt"

import "golang.org/x/sys/unix"

const (
	securebitsNoRoot       = 1 << 0
	securebitsNoRootLocked = 1 << 1
)

// DropPrivileges clears the capabilities gained from entering the user
// namespace. The target must not run with CAP_DAC_OVERRIDE or similar, because
// that would bypass ordinary filesystem permissions inside the sandbox.
func DropPrivileges() error {

	err0 := unix.Prctl(unix.PR_SET_SECUREBITS, uintptr(securebitsNoRoot|securebitsNoRootLocked), 0, 0, 0)

	if err0 == nil {

		for capability := 0; capability <= int(unix.CAP_LAST_CAP); capability++ {
			_ = unix.Prctl(unix.PR_CAPBSET_DROP, uintptr(capability), 0, 0, 0)
		}

		data := [2]unix.CapUserData{}
		err1 := unix.Capset(&unix.CapUserHeader{Version: unix.LINUX_CAPABILITY_VERSION_3}, &data[0])

		if err1 == nil {
			return nil
		} else {
			return fmt.Errorf("landlock.DropPrivileges: Cannot drop capabilities.")
		}

	} else {

		return fmt.Errorf("landlock.DropPrivileges: Cannot set secure bits.")

	}

}
