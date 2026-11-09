package components

import utils_tty "exocomp/utils/tty"
import "strings"

func (renderer *markdownRenderer) addList(indent string, marker string, text string, ordered bool) {

	width := 0

	for index := 0; index < len(indent); index++ {

		if indent[index] == '\t' {
			width += 2
		} else {
			width++
		}

	}

	depth := width / 2
	prefix_text := strings.Repeat("  ", depth) + marker + " "
	marker_kind := "listItem"

	if ordered == true {
		marker_kind = "listEnumeration"
	}

	prefix := []transcriptSpan{
		{Text: prefix_text, Style: markdownStyle(renderer.styles, marker_kind)},
	}
	next := []transcriptSpan{
		{Text: strings.Repeat(" ", utils_tty.StringWidth(prefix_text)), Style: markdownStyle(renderer.styles, "text")},
	}

	spans := renderer.inline(text, "text")
	renderer.appendWrapped(prefix, next, spans, false)

}
