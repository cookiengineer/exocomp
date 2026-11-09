//go:build linux

package landlock

import "errors"
import "fmt"
import "os"
import "runtime"
import "strings"

import "golang.org/x/sys/unix"

// Execute applies landlock confinement and replaces the process with the target.
func Execute(arguments []string) int {

	sandbox, read_paths, target, target_arguments, err0 := parseArguments(arguments)

	if err0 == nil {

		runtime.LockOSThread()

		policy := DefaultPolicy(sandbox, target, read_paths)
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

		} else if errors.Is(err1, errUnsupported) {

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

func parseArguments(arguments []string) (string, []string, string, []string, error) {

	sandbox := ""
	target := ""
	read_paths := make([]string, 0)
	target_arguments := make([]string, 0)
	after_separator := false

	for _, argument := range arguments {

		if after_separator == true {

			if target == "" {
				target = argument
			} else {
				target_arguments = append(target_arguments, argument)
			}

		} else if argument == "--" {
			after_separator = true
		} else if strings.HasPrefix(argument, "--sandbox=") {
			sandbox = strings.TrimSpace(argument[len("--sandbox="):])
		} else if strings.HasPrefix(argument, "--read=") {
			read_paths = append(read_paths, strings.TrimSpace(argument[len("--read="):]))
		} else {
			return "", nil, "", nil, fmt.Errorf("landlock.parseArguments: Invalid argument.")
		}

	}

	if sandbox != "" && target != "" {
		return sandbox, read_paths, target, target_arguments, nil
	} else {
		return "", nil, "", nil, fmt.Errorf("landlock.parseArguments: Invalid argument.")
	}

}
