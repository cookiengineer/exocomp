package components

import "strings"

func firstContentLine(content string) string {

	for _, line := range strings.Split(content, "\n") {

		line = strings.Join(strings.Fields(line), " ")

		if line != "" {
			return line
		}

	}

	return ""

}
