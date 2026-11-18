package findings

import "strings"

func LimitIdentifiers(values []string, maximum int) []string {

	result := make([]string, 0, len(values))

	for _, value := range values {

		value = LimitText(RedactSecrets(strings.TrimSpace(value)), maximum)

		if value != "" {
			result = append(result, value)
		}

	}

	return result

}
