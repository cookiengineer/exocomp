package capabilities

import "encoding/json"
import "strings"

// ParseAnswer parses the JSON answer of a permission dialog.
func ParseAnswer(answer string) ([]Capability, error) {

	answer = strings.TrimSpace(answer)

	if answer == "" {
		return []Capability{}, nil
	} else {

		result := make([]Capability, 0)

		err0 := json.Unmarshal([]byte(answer), &result)

		if err0 == nil {
			return result, nil
		} else {
			return nil, err0
		}

	}

}
