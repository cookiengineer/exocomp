package commands

import "fmt"
import "unicode"

func scanValue(raw string, index *int) (string, error) {

	if *index >= len(raw) {
		return "", nil
	}

	start := *index
	character := raw[*index]

	if character == '"' || character == '\'' {

		token := character
		*index++

		escaped := false

		for *index < len(raw) {

			current := raw[*index]

			if current == '\\' && escaped == false {

				escaped = true
				*index++
				continue

			}

			if current == token && escaped == false {

				*index++
				return raw[start:*index], nil

			}

			escaped = false
			*index++

		}

		return "", fmt.Errorf("commands.Parse: Unterminated string.")

	}

	if character == '{' || character == '[' {

		close_token := byte('}')

		if character == '[' {
			close_token = ']'
		}

		depth := 0
		in_string := false
		escaped := false

		for *index < len(raw) {

			current := raw[*index]
			*index++

			if current == '\\' && in_string == true && escaped == false {

				escaped = true
				continue

			}

			if current == '"' && escaped == false {
				in_string = !in_string
			}

			if in_string == false {

				if current == character {
					depth++
				} else if current == close_token {

					depth--

					if depth == 0 {
						return raw[start:*index], nil
					}

				}

			}

			escaped = false

		}

		return "", fmt.Errorf("commands.Parse: Unterminated nested value.")

	}

	for *index < len(raw) && unicode.IsSpace(rune(raw[*index])) == false {
		*index++
	}

	return raw[start:*index], nil

}
