//go:build windows

package tools

import "os"

// shellKillProcess force-kills the session process on Windows.
func shellKillProcess(pid int) error {

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
