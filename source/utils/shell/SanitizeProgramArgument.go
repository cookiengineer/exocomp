package shell

import utils_capabilities "exocomp/utils/capabilities"
import utils_sandbox "exocomp/utils/sandbox"
import "path/filepath"
import "strings"

// SanitizeProgramArgument applies sandbox path rules and capability grants to a
// program argument. URLs and CIDR lists are network values, not paths, and are
// passed through unchanged.
func SanitizeProgramArgument(sandbox string, capabilities *utils_capabilities.Set, argument string) (string, error) {

	if IsNetworkArgument(argument) == true {
		return argument, nil
	}

	if strings.ContainsAny(argument, "/\\") == true {

		resolved, err0 := utils_sandbox.SanitizePath(sandbox, argument)

		if err0 == nil {

			return resolved, nil

		} else if capabilities != nil && filepath.IsAbs(argument) == true {

			canonical := utils_capabilities.CanonicalPath(argument)

			if capabilities.CanRead(canonical) == true || capabilities.CanWrite(canonical) == true {
				return canonical, nil
			}

		}

		return "", err0

	} else {

		_, err1 := utils_sandbox.CheckPath(sandbox, argument)

		if err1 == nil {
			return argument, nil
		} else {
			return "", err1
		}

	}

}
