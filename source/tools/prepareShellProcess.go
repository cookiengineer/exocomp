package tools

import "os/exec"
import "syscall"

// prepareShellProcess configures the sandbox helper's child process before the
// PTY starts it.
func prepareShellProcess(cmd *exec.Cmd) {

	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	setParentDeathSignal(cmd.SysProcAttr)

}
