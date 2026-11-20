package components

import utils_tty "exocomp/utils/tty"

func writeFooterSpan(buffer *utils_tty.Buffer, x int, y int, text string, style utils_tty.Style, right int) int {

	remaining := right - x

	if remaining <= 0 {
		return 0
	} else {

		text = utils_tty.Truncate(text, remaining)

		if text == "" {
			return 0
		} else {

			buffer.WriteString(x, y, text, style)

			return utils_tty.StringWidth(text)

		}

	}

}
