package shell

import "os/exec"
import "syscall"

// PrepareProcess configures the sandbox helper's child process before the
// PTY starts it.
func PrepareProcess(cmd *exec.Cmd) {

	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	setParentDeathSignal(cmd.SysProcAttr)

}
