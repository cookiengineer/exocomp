package utils

import "strings"

func FormatAgentName(input string) string {

	formatted := make([]byte, 0)

	for i := 0; i < len(input); i++ {

		chr := byte(input[i])

		if chr >= 'A' && chr <= 'Z' {
			formatted = append(formatted, byte(chr))
		} else if chr >= 'a' && chr <= 'z' {
			formatted = append(formatted, byte(chr))
		} else if chr == ' ' {
			formatted = append(formatted, byte(chr))
		}

	}

	tmp1 := strings.TrimSpace(string(formatted))

	if strings.Contains(tmp1, " ") {

		tmp2    := strings.Split(tmp1, " ")
		prename := strings.ToUpper(tmp2[0][0:1]) + strings.ToLower(tmp2[0][1:])
		surname := strings.ToUpper(tmp2[1][0:1]) + strings.ToLower(tmp2[1][1:])

		return prename + " " + surname

	} else {
		return strings.ToUpper(tmp1[0:1]) + strings.ToLower(tmp1[1:])
	}

}
