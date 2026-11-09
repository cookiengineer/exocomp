package commands

import "encoding/json"

func parseSchemaLess(value string) (any, error) {

	if value == "" {
		return "", nil
	}

	if value == "true" {
		return true, nil
	}

	if value == "false" {
		return false, nil
	}

	if value == "null" {
		return nil, nil
	}

	var decoded any
	err := json.Unmarshal([]byte(value), &decoded)

	if err == nil {
		return decoded, nil
	} else {
		return unquote(value), nil
	}

}
