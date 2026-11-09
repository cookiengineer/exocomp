package types

func getConsoleMessageWords(message string) []string {

	result := make([]string, 0)
	chunk := ""

	for m := 0; m < len(message); m++ {

		chr := string(message[m])

		if chr == " " {

			if chunk != "" {
				result = append(result, chunk)
			}

			chunk = ""

		} else {
			chunk += chr
		}

	}

	if chunk != "" {
		result = append(result, chunk)
	}

	return result

}
