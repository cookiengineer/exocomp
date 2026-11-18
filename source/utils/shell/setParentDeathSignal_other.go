//go:build !linux

package shell

import "syscall"

// setParentDeathSignal is a no-op on platforms without a parent-death signal.
func setParentDeathSignal(attributes *syscall.SysProcAttr) {
}
