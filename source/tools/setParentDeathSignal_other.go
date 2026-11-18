//go:build !linux

package tools

import "syscall"

// setParentDeathSignal is a no-op on platforms without a parent-death signal.
func setParentDeathSignal(attributes *syscall.SysProcAttr) {
}
