package tools

import "fmt"

// parseShellInt parses an optional integer argument.
func parseShellInt(arguments map[string]interface{}, name string, fallback int) (int, error) {

	raw, ok := arguments[name]

	if ok == true {

		switch value := raw.(type) {

		case float64:

			if value == float64(int(value)) {
				return int(value), nil
			} else {
				return 0, fmt.Errorf("Invalid parameter \"%s\" is not an integer.", name)
			}

		case int:

			return value, nil

		case int64:

			return int(value), nil

		case uint64:

			return int(value), nil

		default:

			return 0, fmt.Errorf("Invalid parameter \"%s\" is not an integer.", name)

		}

	} else {
		return fallback, nil
	}

}
