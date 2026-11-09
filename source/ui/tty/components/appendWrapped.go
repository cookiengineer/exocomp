package components

import utils_tty "exocomp/utils/tty"

func (renderer *markdownRenderer) appendWrapped(first []transcriptSpan, next []transcriptSpan, spans []transcriptSpan, fill bool) {

	first_width := 0

	for _, span := range first {
		first_width += utils_tty.StringWidth(span.Text)
	}

	next_width := 0

	for _, span := range next {
		next_width += utils_tty.StringWidth(span.Text)
	}

	prefix_width := max(first_width, next_width)
	content_width := max(1, renderer.width-prefix_width)
	wrapped := wrapSpans(spans, content_width)

	if len(wrapped) == 0 {
		wrapped = [][]transcriptSpan{{}}
	}

	for index, line := range wrapped {

		prefix := next

		if index == 0 {
			prefix = first
		}

		combined := make([]transcriptSpan, 0)
		combined = append(combined, prefix...)
		combined = append(combined, line...)

		result := transcriptLine{
			Spans: compactSpans(combined),
		}

		if fill == true {
			result.HasFill = true
			result.Fill = markdownStyle(renderer.styles, "codeBlock")
		}

		renderer.lines = append(renderer.lines, result)

	}

}
