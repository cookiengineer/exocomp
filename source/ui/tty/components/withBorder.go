package components

import utils_tty "exocomp/utils/tty"

func withBorder(lines []transcriptLine, style utils_tty.Style) []transcriptLine {

	for index := range lines {

		lines[index].HasBorder = true
		lines[index].Border = style

	}

	return lines

}
