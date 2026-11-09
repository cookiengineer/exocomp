package tty

import "strings"

func splitRunes(text string, width int) []string {

	result := make([]string, 0)
	chunk := strings.Builder{}
	used := 0

	for _, value := range text {

		value_width := RuneWidth(value)

		if used+value_width > width && used > 0 {

			result = append(result, chunk.String())
			chunk.Reset()
			used = 0

		}

		chunk.WriteRune(value)
		used += value_width

	}

	if chunk.Len() > 0 {
		result = append(result, chunk.String())
	}

	if len(result) == 0 {
		result = append(result, "")
	}

	return result

}
