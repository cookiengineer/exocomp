package tty

import "strings"

func Truncate(text string, width int) string {

	if width <= 0 {
		return ""
	} else if StringWidth(text) <= width {
		return text
	} else {

		builder := strings.Builder{}
		used := 0

		for _, value := range text {

			value_width := RuneWidth(value)

			if used+value_width > width {
				break
			} else {

				builder.WriteRune(value)
				used += value_width

			}
		}

		return builder.String()

	}
}
