package assets

import "encoding/json"
import "fmt"
import "strings"

func decodeImportLines(data []byte) ([]json.RawMessage, int, error) {

	trimmed := strings.TrimSpace(string(data))

	if trimmed == "" {
		return make([]json.RawMessage, 0), 0, nil
	} else {

		if strings.HasPrefix(trimmed, "[") == true {

			result := make([]json.RawMessage, 0)

			if err := json.Unmarshal([]byte(trimmed), &result); err == nil {
				return result, 0, nil
			} else {
				return nil, 0, fmt.Errorf("Cannot parse results.")
			}
		}

		result := make([]json.RawMessage, 0)
		invalid := 0

		for _, line := range strings.Split(trimmed, "\n") {

			line = strings.TrimSpace(line)

			if line == "" || strings.HasPrefix(line, "#") == true {
				continue
			} else {

				if json.Valid([]byte(line)) == false {
					invalid++
					continue
				}

				result = append(result, json.RawMessage(line))

			}
		}

		return result, invalid, nil

	}
}
