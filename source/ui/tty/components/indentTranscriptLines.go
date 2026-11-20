package components

import utils_tty "exocomp/utils/tty"

func indentTranscriptLines(lines []transcriptLine, prefix string, style utils_tty.Style) []transcriptLine {

	result := make([]transcriptLine, 0, len(lines))

	for _, line := range lines {

		if len(line.Spans) == 0 && line.HasFill == false {

			result = append(result, line)
			continue

		}

		prefix_style := style

		if line.HasFill == true {
			prefix_style = line.Fill
		}

		spans := make([]transcriptSpan, 0, len(line.Spans)+1)
		spans = append(spans, transcriptSpan{Text: prefix, Style: prefix_style})
		spans = append(spans, line.Spans...)
		line.Spans = compactSpans(spans)
		result = append(result, line)

	}

	return result

}
