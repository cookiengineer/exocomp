//go:build windows

package shell

import "os"

// KillProcess force-kills the session process on Windows.
func KillProcess(pid int) error {

	if pid > 0 {
		process, err := os.FindProcess(pid)
		if err == nil {
			return process.Kill()
		} else {
			return err
		}
	} else {
		return nil
	}

}
