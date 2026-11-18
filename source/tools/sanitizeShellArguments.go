package tools

import utils_capabilities "exocomp/utils/capabilities"

// sanitizeShellArguments applies the same path rules as programs.Execute to
// arguments that contain a path separator.
func sanitizeShellArguments(sandbox string, capabilities *utils_capabilities.Set, arguments []string) ([]string, error) {

	result := make([]string, 0, len(arguments))

	for _, argument := range arguments {

		resolved, err0 := sanitizeProgramArgument(sandbox, capabilities, argument)

		if err0 == nil {
			result = append(result, resolved)
		} else {
			return nil, err0
		}

	}

	return result, nil

}
