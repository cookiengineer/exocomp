package types

func getConsoleIndent(offset int) string {

	result := ""

	if offset > 0 {

		result = "|"

		for o := 1; o < offset; o++ {
			result += "|"
		}

	}

	return result

}

