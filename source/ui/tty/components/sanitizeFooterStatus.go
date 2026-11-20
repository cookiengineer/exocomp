package components

import "strings"

func sanitizeFooterStatus(value string) string {

	builder := strings.Builder{}

	for _, character := range value {

		if character == '\n' || character == '\r' || character == '\t' {
			builder.WriteRune(' ')
		} else if character >= 32 && character != 127 {
			builder.WriteRune(character)
		}

	}

	return strings.Join(strings.Fields(builder.String()), " ")

}
