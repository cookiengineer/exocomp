package components

import "strings"

func renderMarkdownPlain(text string, width int) []transcriptLine {

	renderer := newMarkdownRenderer(width)

	for _, line := range strings.Split(text, "\n") {

		if line == "" {
			renderer.lines = append(renderer.lines, transcriptLine{})
		} else {

			spans := []transcriptSpan{
				{Text: line, Style: markdownStyle(renderer.styles, "text")},
			}
			renderer.appendWrapped(nil, nil, spans, false)

		}

	}

	return renderer.lines

}
