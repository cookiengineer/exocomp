package tty

import "strings"

func Truncate(text string, width int) string {

	if width <= 0 {
		return ""
	}

	if StringWidth(text) <= width {
		return text
	}

	builder := strings.Builder{}
	used := 0

	for _, value := range text {

		value_width := RuneWidth(value)

		if used+value_width > width {
			break
		}

		builder.WriteRune(value)
		used += value_width

	}

	return builder.String()

}
