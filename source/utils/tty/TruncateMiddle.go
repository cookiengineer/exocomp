package tty

func TruncateMiddle(text string, width int) string {

	if width <= 0 {
		return ""
	}

	if StringWidth(text) <= width {
		return text
	}

	if width <= 3 {
		return Truncate(text, width)
	}

	runes := []rune(text)
	left := make([]rune, 0)
	right := make([]rune, 0)
	left_width := 0
	right_width := 0
	budget := width - 1

	for l := 0; l < len(runes); l++ {

		value_width := RuneWidth(runes[l])

		if left_width+value_width > budget/2 {
			break
		}

		left = append(left, runes[l])
		left_width += value_width

	}

	for r := len(runes) - 1; r >= 0; r-- {

		value_width := RuneWidth(runes[r])

		if right_width+value_width > budget-left_width {
			break
		}

		right = append([]rune{runes[r]}, right...)
		right_width += value_width

	}

	return string(left) + "…" + string(right)

}
