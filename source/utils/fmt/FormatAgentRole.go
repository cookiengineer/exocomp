package fmt

import "strings"

func FormatAgentRole(input string) string {

	formatted := make([]byte, 0)

	for i := 0; i < len(input); i++ {

		chr := byte(input[i])

		if chr >= 'A' && chr <= 'Z' {
			formatted = append(formatted, byte(chr))
		} else if chr >= 'a' && chr <= 'z' {
			formatted = append(formatted, byte(chr))
		}

	}

	tmp1 := strings.TrimSpace(string(formatted))

	return strings.ToLower(tmp1)

}
