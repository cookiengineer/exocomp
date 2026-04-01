package utils

import "strings"

func FormatMultiLine(input string) string {

	formatted := make([]byte, 0)

	for i := 0; i < len(input); i++ {

		chr := byte(input[i])

		// isPrintableASCII also filters \t, \r and \n
		if isPrintableASCII(chr) {
			formatted = append(formatted, chr)
		} else if chr == '\t' || chr == '\n' {
			formatted = append(formatted, chr)
		}

	}

	return strings.TrimSpace(string(formatted))

}
