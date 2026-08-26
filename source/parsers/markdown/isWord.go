package markdown

func isWord(raw string) bool {

	result := true

	for _, chr := range raw {

		if chr >= 'a' && chr <= 'z' {
			// Do Nothing
		} else if chr >= 'A' && chr <= 'Z' {
			// Do Nothing
		} else {
			result = false
			break
		}

	}

	return result

}
