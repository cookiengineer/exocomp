package utils

import "strings"

func FormatFilePath(input string) string {

	formatted := make([]byte, 0)

	for i := 0; i < len(input); i++ {

		chr := byte(input[i])

		if chr >= '0' && chr <= '9' {
			formatted = append(formatted, byte(chr))
		} else if chr >= 'A' && chr <= 'Z' {
			formatted = append(formatted, byte(chr))
		} else if chr >= 'a' && chr <= 'z' {
			formatted = append(formatted, byte(chr))
		} else if chr == '/' || chr == '\\' {
			formatted = append(formatted, byte(chr))
		} else if chr == '.' || chr == '_' || chr == '-' {
			formatted = append(formatted, byte(chr))
		} else if chr == ' ' {
			formatted = append(formatted, byte(chr))
		}

	}

	return strings.TrimSpace(string(formatted))

}
