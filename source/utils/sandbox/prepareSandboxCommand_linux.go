//go:build linux

package sandbox

import "os"
import "os/exec"
import "syscall"

// prepareSandboxCommand puts the sandbox helper into a private user, mount and
// (when network access is denied) network namespace. The helper uses the mount
// namespace to mask the private .exocomp state folder before applying landlock
// and executing the target. The private network namespace is the only network
// boundary on Linux: landlock cannot match individual hosts, so a granted
// network.Connect capability keeps the helper in the host network namespace.
func prepareSandboxCommand(cmd *exec.Cmd, has_network bool) {

	clone_flags := uintptr(syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS)

	if has_network == false {
		clone_flags = clone_flags | uintptr(syscall.CLONE_NEWNET)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: clone_flags,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
		GidMappingsEnableSetgroups: false,
	}

}
