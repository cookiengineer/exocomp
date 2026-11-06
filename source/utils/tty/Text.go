package tty

import "strings"
import "unicode"

func RuneWidth(value rune) int {

	if value == 0 {
		return 0
	}

	if value < 32 || (value >= 0x7f && value < 0xa0) {
		return 0
	}

	if unicode.Is(unicode.Mn, value) || unicode.Is(unicode.Me, value) || unicode.Is(unicode.Cf, value) {
		return 0
	}

	if isWideRune(value) {
		return 2
	}

	return 1

}

func isWideRune(value rune) bool {

	return (value >= 0x1100 && value <= 0x115f) ||
		value == 0x2329 ||
		value == 0x232a ||
		(value >= 0x2e80 && value <= 0xa4cf && value != 0x303f) ||
		(value >= 0xac00 && value <= 0xd7a3) ||
		(value >= 0xf900 && value <= 0xfaff) ||
		(value >= 0xfe10 && value <= 0xfe19) ||
		(value >= 0xfe30 && value <= 0xfe6f) ||
		(value >= 0xff00 && value <= 0xff60) ||
		(value >= 0xffe0 && value <= 0xffe6) ||
		(value >= 0x1f300 && value <= 0x1f64f) ||
		(value >= 0x1f900 && value <= 0x1f9ff) ||
		(value >= 0x20000 && value <= 0x3fffd)

}

func StringWidth(text string) int {

	width := 0

	for _, value := range text {
		width += RuneWidth(value)
	}

	return width

}

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

func TruncateRight(text string, width int) string {

	if StringWidth(text) <= width {
		return text
	}

	return Truncate(text, max(0, width-1)) + "…"

}

func PadRight(text string, width int) string {

	length := StringWidth(text)

	if length >= width {
		return text
	}

	return text + strings.Repeat(" ", width-length)

}

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

func splitRunes(text string, width int) []string {

	result := make([]string, 0)
	chunk := strings.Builder{}
	used := 0

	for _, value := range text {

		value_width := RuneWidth(value)

		if used+value_width > width && used > 0 {
			result = append(result, chunk.String())
			chunk.Reset()
			used = 0
		}

		chunk.WriteRune(value)
		used += value_width

	}

	if chunk.Len() > 0 {
		result = append(result, chunk.String())
	}

	if len(result) == 0 {
		result = append(result, "")
	}

	return result

}
