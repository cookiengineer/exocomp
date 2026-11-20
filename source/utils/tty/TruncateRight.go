package tty

func TruncateRight(text string, width int) string {

	if StringWidth(text) <= width {
		return text
	} else {
		return Truncate(text, max(0, width-1)) + "…"
	}
}
