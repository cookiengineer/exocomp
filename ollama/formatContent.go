package ollama

import "strings"

func formatContent(message string) []string {

	result := make([]string, 0)
	lines  := strings.Split(message, "\n")

	for l := 0; l < len(lines); l++ {

		line := lines[l]

		if strings.HasPrefix(line, "###") {
		} else {
		}

		result = append(result, line)

	}

	return result

}
