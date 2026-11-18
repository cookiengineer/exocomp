package tools

import "strings"

// sanitizeShellArguments applies the same path rules as programs.Execute to
// arguments that contain a path separator.
func sanitizeShellArguments(sandbox string, arguments []string) ([]string, error) {

	result := make([]string, 0, len(arguments))

	for _, argument := range arguments {

		if strings.ContainsAny(argument, "/\\") == true {

			resolved, err0 := sanitizeSandboxPath(sandbox, argument)

			if err0 == nil {
				result = append(result, resolved)
			} else {
				return nil, err0
			}

		} else {

			_, err1 := checkSandboxPath(sandbox, argument)

			if err1 == nil {
				result = append(result, argument)
			} else {
				return nil, err1
			}

		}

	}

	return result, nil

}
