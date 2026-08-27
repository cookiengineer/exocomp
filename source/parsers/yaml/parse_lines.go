package yaml

import "strings"

func parse_lines(input string) []parser_line {

	lines  := strings.Split(input, "\n")
	result := make([]parser_line, 0)

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

			result  = append(result, parser_line{
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

