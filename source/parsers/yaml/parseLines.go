package yaml

import "strings"

func parseLines(input string) []parserLine {

	lines := strings.Split(input, "\n")
	result := make([]parserLine, 0)

	for index, line := range lines {

		tmp := strings.TrimRight(line, " ")

		if strings.TrimSpace(tmp) != "" {

			indent := int(0)

			for _, character := range tmp {

				if character == ' ' || character == '\t' {

					indent++
					continue

				} else {
					break
				}

			}

			result = append(result, parserLine{
				Number: index + 1,
				Indent: indent,
				Text:   strings.TrimSpace(tmp),
			})

		} else {
			continue
		}

	}

	return result

}
