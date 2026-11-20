//go:build !linux

package sandbox

import "os/exec"

// prepareSandboxCommand is a no-op on platforms without Linux namespaces.
func prepareSandboxCommand(cmd *exec.Cmd, has_network bool) {
}
