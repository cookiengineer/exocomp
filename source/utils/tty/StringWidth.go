package tty

func StringWidth(text string) int {

	width := 0

	for _, value := range text {
		width += RuneWidth(value)
	}

	return width

}
