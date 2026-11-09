package components

import utils_tty "exocomp/utils/tty"
import "strings"

func wrapSpans(spans []transcriptSpan, width int) [][]transcriptSpan {

	width = max(1, width)

	runes := make([]styledRune, 0)

	for _, span := range spans {

		for _, value := range span.Text {

			if value == '\t' {

				for tab := 0; tab < 4; tab++ {
					runes = append(runes, styledRune{Rune: ' ', Style: span.Style})
				}

			} else {
				runes = append(runes, styledRune{Rune: value, Style: span.Style})
			}

		}

	}

	lines := make([][]styledRune, 0)
	current := make([]styledRune, 0)
	current_width := 0

	flush := func() {
		lines = append(lines, current)
		current = make([]styledRune, 0)
		current_width = 0
	}

	for _, item := range runes {

		if item.Rune == '\r' {
			continue
		}

		if item.Rune == '\n' {
			flush()
			continue
		}

		item_width := utils_tty.RuneWidth(item.Rune)

		if item_width <= 0 {

			if len(current) > 0 {
				current = append(current, item)
			}

			continue

		}

		if current_width+item_width > width && len(current) > 0 {

			break_at := -1

			for index := len(current) - 1; index >= 0; index-- {

				if current[index].Rune == ' ' {
					break_at = index
					break
				}

			}

			if break_at > 0 {

				line := append([]styledRune{}, current[0:break_at]...)
				lines = append(lines, line)
				remainder := current[break_at+1:]

				for len(remainder) > 0 && remainder[0].Rune == ' ' {
					remainder = remainder[1:]
				}

				current = append([]styledRune{}, remainder...)
				current_width = 0

				for _, value := range current {
					current_width += utils_tty.RuneWidth(value.Rune)
				}

			} else {

				lines = append(lines, current)
				current = make([]styledRune, 0)
				current_width = 0

			}

		}

		current = append(current, item)
		current_width += item_width

	}

	if len(current) > 0 || len(lines) == 0 {
		lines = append(lines, current)
	}

	if len(lines) > 1 && len(lines[len(lines)-1]) == 0 {
		lines = lines[0 : len(lines)-1]
	}

	result := make([][]transcriptSpan, 0)

	for _, line := range lines {

		line_result := make([]transcriptSpan, 0)
		builder := strings.Builder{}
		current_style := utils_tty.Style{}
		has_style := false

		for _, item := range line {

			if has_style == false {
				current_style = item.Style
				has_style = true
			}

			if item.Style.Equals(current_style) == false {
				line_result = append(line_result, transcriptSpan{Text: builder.String(), Style: current_style})
				builder.Reset()
				current_style = item.Style
			}

			builder.WriteRune(item.Rune)

		}

		if builder.Len() > 0 {
			line_result = append(line_result, transcriptSpan{Text: builder.String(), Style: current_style})
		}

		result = append(result, line_result)

	}

	return result

}
