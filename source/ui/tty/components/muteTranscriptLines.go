package components

import "exocomp/ui/tty/colors"

func muteTranscriptLines(lines []transcriptLine) []transcriptLine {

	for index := range lines {

		if lines[index].HasFill == true {
			lines[index].Fill = lines[index].Fill.WithForeground(colors.TextMuted)
		}

		for span := range lines[index].Spans {
			lines[index].Spans[span].Style = lines[index].Spans[span].Style.WithForeground(colors.TextMuted)
		}

	}

	return lines

}
