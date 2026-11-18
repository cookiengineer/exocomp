package findings

import "strings"
import "unicode/utf8"

func LimitText(value string, maximum int) string {

	value = strings.TrimSpace(value)

	if len(value) > maximum {

		value = value[:maximum]

		for len(value) > 0 && utf8.ValidString(value) == false {
			value = value[:len(value)-1]
		}

		return strings.TrimSpace(value) + "\n[TRUNCATED]"

	} else {
		return value
	}

}
