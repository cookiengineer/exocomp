package fmt

import "strings"

func FormatAgentModel(input string) string {

	formatted := make([]byte, 0)

	for i := 0; i < len(input); i++ {

		chr := byte(input[i])

		if chr >= 'a' && chr <= 'z' {
			formatted = append(formatted, byte(chr))
		} else if chr == '-' || chr == ':' {
			formatted = append(formatted, byte(chr))
		}

	}

	return strings.ToLower(string(formatted))

}
