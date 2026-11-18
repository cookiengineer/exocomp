//go:build linux

package sandbox

import "os"
import "os/exec"
import "syscall"

// prepareSandboxCommand puts the sandbox helper into a private user and mount
// namespace. The helper uses the mount namespace to mask the private .exocomp
// state folder before applying landlock and executing the target.
func prepareSandboxCommand(cmd *exec.Cmd) {

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
		GidMappingsEnableSetgroups: false,
	}

}
