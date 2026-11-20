package capabilities

import "fmt"

// Parse parses the capabilities array of humans.AskPermission.
func Parse(raw any) ([]Capability, error) {

	result := make([]Capability, 0)

	switch values := raw.(type) {

	case []Capability:

		return append([]Capability{}, values...), nil

	case []any:

		for index, value := range values {

			entry, ok := value.(map[string]any)

			if ok == true {

				capability_type, ok1 := entry["type"].(string)
				scope, ok2 := entry["scope"].(string)

				if ok1 == true && ok2 == true {
					result = append(result, Capability{Type: capability_type, Scope: scope})
				} else {
					return nil, fmt.Errorf("capabilities.Parse: Invalid parameter \"capabilities[%d]\" must contain string \"type\" and \"scope\".", index)
				}

			} else {
				return nil, fmt.Errorf("capabilities.Parse: Invalid parameter \"capabilities[%d]\" is not an object.", index)
			}

		}

		return result, nil

	default:
		return nil, fmt.Errorf("capabilities.Parse: Invalid parameter \"capabilities\" is not an array of capabilities.")

	}

}
