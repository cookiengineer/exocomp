//go:build linux

package shell

import "syscall"

// setParentDeathSignal makes the sandbox helper and its target die when the
// agent process dies, so interactive sessions cannot be orphaned on Linux.
func setParentDeathSignal(attributes *syscall.SysProcAttr) {
	attributes.Pdeathsig = syscall.SIGKILL
}
