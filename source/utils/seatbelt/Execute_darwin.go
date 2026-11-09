//go:build darwin

package seatbelt

import "fmt"
import "os"
import "path/filepath"
import "syscall"

// Execute applies seatbelt confinement and replaces the process with the target.
func Execute(arguments []string) int {

	sandbox, read_paths, target, target_arguments, err0 := parseArguments(arguments)

	if err0 == nil {

		sandbox_abs, err1 := filepath.Abs(sandbox)

		if err1 == nil {

			target_abs, err2 := filepath.Abs(target)

			if err2 == nil {

				sandbox_exec, err3 := sandboxExecutable()

				if err3 == nil {

					policy := DefaultPolicy(sandbox_abs, target_abs, read_paths)
					profile, err4 := BuildProfile(policy)

					if err4 == nil {

						argv := []string{
							sandbox_exec,
							"-p",
							profile.String(),
							target_abs,
						}

						argv = append(argv, target_arguments...)

						err5 := syscall.Exec(sandbox_exec, argv, os.Environ())

						if err5 == nil {
							return 0
						} else {
							fmt.Fprint(os.Stderr, "seatbelt.Execute: Cannot execute program.\n")
							return 126
						}

					} else {
						fmt.Fprint(os.Stderr, "seatbelt.Execute: Cannot apply sandbox confinement.\n")
						return 125
					}

				} else {
					fmt.Fprint(os.Stderr, "seatbelt.Execute: Sandboxing is not supported on this system.\n")
					return 125
				}

			} else {
				fmt.Fprint(os.Stderr, "seatbelt.Execute: Invalid arguments.\n")
				return 2
			}

		} else {
			fmt.Fprint(os.Stderr, "seatbelt.Execute: Invalid arguments.\n")
			return 2
		}

	} else {
		fmt.Fprint(os.Stderr, "seatbelt.Execute: Invalid arguments.\n")
		return 2
	}

}
