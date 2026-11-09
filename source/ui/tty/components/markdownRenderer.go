package components

import utils_tty "exocomp/utils/tty"
import "strings"

type markdownRenderer struct {
	width  int
	styles map[string]utils_tty.Style
	lines  []transcriptLine
}

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

func (renderer *markdownRenderer) inline(text string, default_kind string) []transcriptSpan {

	result := make([]transcriptSpan, 0)
	default_style := markdownStyle(renderer.styles, default_kind)
	builder := strings.Builder{}

	flush := func() {

		if builder.Len() > 0 {

			result = append(result, transcriptSpan{Text: builder.String(), Style: default_style})
			builder.Reset()

		}

	}

	for index := 0; index < len(text); {

		if text[index] == '`' {

			end := strings.Index(text[index+1:], "`")

			if end >= 0 {

				flush()
				inner := text[index+1 : index+1+end]
				result = append(result, transcriptSpan{Text: inner, Style: markdownStyle(renderer.styles, "code")})
				index = index + end + 2
				continue

			}

		}

		if strings.HasPrefix(text[index:], "**") == true {

			end := strings.Index(text[index+2:], "**")

			if end >= 0 {

				flush()
				inner := text[index+2 : index+2+end]
				result = append(result, transcriptSpan{Text: inner, Style: markdownStyle(renderer.styles, "strong")})
				index = index + end + 4
				continue

			}

		}

		if strings.HasPrefix(text[index:], "~~") == true {

			end := strings.Index(text[index+2:], "~~")

			if end >= 0 {

				flush()
				inner := text[index+2 : index+2+end]
				result = append(result, transcriptSpan{Text: inner, Style: markdownStyle(renderer.styles, "strike")})
				index = index + end + 4
				continue

			}

		}

		if text[index] == '*' {

			end := strings.Index(text[index+1:], "*")

			if end > 0 {

				flush()
				inner := text[index+1 : index+1+end]
				result = append(result, transcriptSpan{Text: inner, Style: markdownStyle(renderer.styles, "emphasis")})
				index = index + end + 2
				continue

			}

		}

		if strings.HasPrefix(text[index:], "![") == true {

			label, url, end, ok := parseMarkdownLink(text, index, true)

			if ok == true {

				flush()
				result = append(result, transcriptSpan{Text: label, Style: markdownStyle(renderer.styles, "imageText")})

				if url != "" {
					result = append(result, transcriptSpan{Text: " (" + url + ")", Style: markdownStyle(renderer.styles, "linkUrl")})
				}

				index = end
				continue

			}

		}

		if text[index] == '[' {

			label, url, end, ok := parseMarkdownLink(text, index, false)

			if ok == true {

				flush()
				result = append(result, transcriptSpan{Text: label, Style: markdownStyle(renderer.styles, "link")})

				if url != "" && url != label {
					result = append(result, transcriptSpan{Text: " (" + url + ")", Style: markdownStyle(renderer.styles, "linkUrl")})
				}

				index = end
				continue

			}

		}

		builder.WriteByte(text[index])
		index++

	}

	flush()

	return result

}
