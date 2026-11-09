package tty

import "strings"

func Wrap(text string, width int) []string {

	if width <= 0 {
		return []string{""}
	}

	result := make([]string, 0)
	paragraphs := strings.Split(text, "\n")

	for _, paragraph := range paragraphs {

		if paragraph == "" {

			result = append(result, "")
			continue

		}

		line := strings.Builder{}
		line_width := 0

		for _, word := range strings.Split(paragraph, " ") {

			word_width := StringWidth(word)

			if line_width == 0 {

				if word_width > width {

					chunks := splitRunes(word, width)
					result = append(result, chunks[0:len(chunks)-1]...)
					line.WriteString(chunks[len(chunks)-1])
					line_width = StringWidth(chunks[len(chunks)-1])

				} else {

					line.WriteString(word)
					line_width = word_width

				}

			} else if line_width+1+word_width <= width {

				line.WriteString(" ")
				line.WriteString(word)
				line_width += 1 + word_width

			} else {

				result = append(result, line.String())
				line.Reset()
				line.WriteString(word)
				line_width = word_width

			}

		}

		if line.Len() > 0 {
			result = append(result, line.String())
		}

	}

	if len(result) == 0 {
		result = append(result, "")
	}

	return result

}
