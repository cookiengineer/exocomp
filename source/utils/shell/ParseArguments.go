package shell

import "fmt"

// ParseArguments parses the optional string array argument of shell.Open.
func ParseArguments(raw any) ([]string, error) {

	result := make([]string, 0)

	if raw == nil {

		return result, nil

	} else {

		switch values := raw.(type) {

		case []string:

			return append([]string{}, values...), nil

		case []interface{}:

			valid := true

			for _, value := range values {

				tmp, ok := value.(string)

				if ok == true {
					result = append(result, tmp)
				} else {
					valid = false
					break
				}

			}

			if valid == true {
				return result, nil
			} else {
				return nil, fmt.Errorf("Invalid parameter \"arguments\" is not an array of strings.")
			}

		default:

			return nil, fmt.Errorf("Invalid parameter \"arguments\" is not an array of strings.")

		}

	}

}
