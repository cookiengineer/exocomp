//go:build linux

package landlock

import "errors"
import "fmt"
import "os"
import "runtime"

import "golang.org/x/sys/unix"

// Execute applies landlock confinement and replaces the process with the target.
func Execute(arguments []string) int {

	sandbox, read_paths, write_paths, allow_network, target, target_arguments, err0 := parseArguments(arguments)

	if err0 == nil {

		runtime.LockOSThread()

		// NOTE: The parent puts the helper into a private network namespace
		// when network access is denied. Bring its loopback interface up so
		// local test servers and loopback clients keep working without
		// granting access to the host network.
		if allow_network == false {

			err_loopback := raiseLoopback()

			if err_loopback != nil {

				fmt.Fprint(os.Stderr, "landlock.Execute: Cannot isolate network namespace.\n")
				return 125

			}

		}

		err_mask := maskPrivateState(sandbox)

		if err_mask != nil {

			fmt.Fprint(os.Stderr, "landlock.Execute: Cannot hide private state folder.\n")
			return 125

		}

		policy := DefaultPolicy(sandbox, target, read_paths, write_paths)
		err1 := restrictPaths(policy)

		if err1 == nil {

			argv := append([]string{target}, target_arguments...)
			err2 := unix.Exec(target, argv, os.Environ())

			if err2 == nil {
				return 0
			} else {

				fmt.Fprint(os.Stderr, "landlock.Execute: Cannot execute program.\n")
				return 126

			}

		} else if errors.Is(err1, errUnsupported) == true {

			fmt.Fprint(os.Stderr, "landlock.Execute: Sandboxing is not supported on this system.\n")
			return 125

		} else {

			fmt.Fprint(os.Stderr, "landlock.Execute: Cannot apply sandbox confinement.\n")
			return 125

		}

	} else {

		fmt.Fprint(os.Stderr, "landlock.Execute: Invalid arguments.\n")
		return 2

	}

}
