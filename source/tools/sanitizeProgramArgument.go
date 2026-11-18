package tools

import utils_capabilities "exocomp/utils/capabilities"
import "path/filepath"
import "strings"

// sanitizeProgramArgument applies sandbox path rules and capability grants to a
// program argument.
func sanitizeProgramArgument(sandbox string, capabilities *utils_capabilities.Set, argument string) (string, error) {

	if strings.ContainsAny(argument, "/\\") == true {

		resolved, err0 := sanitizeSandboxPath(sandbox, argument)

		if err0 == nil {

			return resolved, nil

		} else if capabilities != nil && filepath.IsAbs(argument) == true {

			canonical, err1 := filepath.EvalSymlinks(argument)

			if err1 != nil {

				parent, err2 := filepath.EvalSymlinks(filepath.Dir(argument))

				if err2 == nil {
					canonical = filepath.Join(parent, filepath.Base(argument))
				}

			}

			if canonical != "" {

				if capabilities.CanRead(canonical) == true || capabilities.CanWrite(canonical) == true {
					return canonical, nil
				}

			}

		}

		return "", err0

	} else {

		_, err3 := checkSandboxPath(sandbox, argument)

		if err3 == nil {
			return argument, nil
		} else {
			return "", err3
		}

	}

}
