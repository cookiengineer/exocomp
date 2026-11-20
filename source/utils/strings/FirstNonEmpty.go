package strings

import "strings"

// FirstNonEmpty returns the first trimmed value that is not empty.
func FirstNonEmpty(values ...string) string {

	for _, value := range values {

		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}

	}

	return ""

}
