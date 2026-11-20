package commands

import "encoding/json"

func parseSchemaLess(value string) (any, error) {

	if value == "" {
		return "", nil
	} else if value == "true" {
		return true, nil
	} else if value == "false" {
		return false, nil
	} else if value == "null" {
		return nil, nil
	} else {

		var decoded any
		err := json.Unmarshal([]byte(value), &decoded)

		if err == nil {
			return decoded, nil
		} else {
			return unquote(value), nil
		}

	}
}
