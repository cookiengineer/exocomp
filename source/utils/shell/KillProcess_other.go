//go:build !windows

package shell

import "syscall"

// KillProcess force-kills the session's process group.
func KillProcess(pid int) error {

	if pid > 0 {

		err := syscall.Kill(-pid, syscall.SIGKILL)
		if err == syscall.ESRCH {
			return nil
		} else {
			return err
		}

	} else {
		return nil
	}

}
