package strings

func SeekNested(raw string, start int, start_token byte, close_token byte) []byte {

	result := []byte{start_token}
	depth := 0
	in_string := false
	escaped := false

	for index := start + 1; index < len(raw); index++ {

		result = append(result, byte(raw[index]))

		chr := raw[index]

		if chr == '\\' && in_string && escaped == false {
			escaped = true
			continue
		}

		if chr == '"' && escaped == false {
			in_string = !in_string
		}

		if in_string == false {

			if chr == start_token {

				depth++

			} else if chr == close_token {

				depth--

				if depth == 0 {
					break
				}

			}

		}

		escaped = false

	}

	return result

}


