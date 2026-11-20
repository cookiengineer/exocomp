package tty

import "strings"

func PadRight(text string, width int) string {

	length := StringWidth(text)

	if length >= width {
		return text
	} else {
		return text + strings.Repeat(" ", width-length)
	}
}
