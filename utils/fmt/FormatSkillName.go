package utils

import "strings"

func FormatSkillName(input string) string {

	input = strings.ToLower(input)

	formatted     := make([]byte, 0)
	last_was_dash := false

	for i := 0; i < len(input); i++ {

		chr := byte(input[i])

		if chr >= 'a' && chr <= 'z' {

			formatted = append(formatted, byte(chr))
			last_was_dash = false

		} else if chr == '-' {

			if i >= 1 && i < len(input) - 1 {

				if last_was_dash == false {
					formatted = append(formatted, byte(chr))
					last_was_dash = true
				}

			}

		}

	}

	return string(formatted)

}

