package strings

func Seek(raw string, start int, token byte) []byte {

	result := []byte{token}
	escaped := false

	for index := start + 1; index < len(raw); index++ {

		result = append(result, byte(raw[index]))

		chr := raw[index]

		if chr == '\\' && escaped == false {

			escaped = true
			continue

		}

		if chr == token && escaped == false {
			break
		}

		escaped = false

	}

	return result

}
