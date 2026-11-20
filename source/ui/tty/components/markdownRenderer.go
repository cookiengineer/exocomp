package components

import parser_markdown "exocomp/parsers/markdown"
import "strconv"
import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"
import "strings"

type markdownRenderer struct {
	width  int
	styles map[string]utils_tty.Style
	lines  []transcriptLine
}

func newMarkdownRenderer(width int) *markdownRenderer {

	renderer := &markdownRenderer{
		width:  max(1, width),
		styles: make(map[string]utils_tty.Style),
		lines:  make([]transcriptLine, 0),
	}

	renderer.styles["text"] = styles.MarkdownText
	renderer.styles["heading"] = styles.MarkdownHeading
	renderer.styles["strong"] = styles.MarkdownStrong
	renderer.styles["emphasis"] = styles.MarkdownEmphasis
	renderer.styles["code"] = styles.MarkdownCode
	renderer.styles["codeBlock"] = styles.MarkdownCodeBlock
	renderer.styles["link"] = styles.MarkdownLink
	renderer.styles["linkUrl"] = styles.MarkdownLinkURL
	renderer.styles["blockQuote"] = styles.MarkdownBlockQuote
	renderer.styles["listItem"] = styles.MarkdownListItem
	renderer.styles["listEnumeration"] = styles.MarkdownListEnumeration
	renderer.styles["horizontalRule"] = styles.MarkdownHorizontalRule
	renderer.styles["image"] = styles.MarkdownImage
	renderer.styles["imageText"] = styles.MarkdownImageText
	renderer.styles["strike"] = styles.MarkdownStrike

	return renderer

}

func (renderer *markdownRenderer) addList(indent string, marker string, children []*parser_markdown.Element, ordered bool) {

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

	spans := renderer.inlineChildren(children, "text")
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

func (renderer *markdownRenderer) inlineChildren(children []*parser_markdown.Element, default_kind string) []transcriptSpan {

	result := make([]transcriptSpan, 0)
	default_style := markdownStyle(renderer.styles, default_kind)
	last_line := 0

	for _, child := range children {

		if child == nil {
			continue
		} else {

			spans := renderer.inlineElement(child, default_kind)

			if len(spans) == 0 {
				continue
			} else {

				if last_line > 0 && child.Line > last_line {
					result = append(result, transcriptSpan{Text: "\n", Style: default_style})
				} else if len(result) > 0 && needsMarkdownSpace(result, spans) == true {
					result = append(result, transcriptSpan{Text: " ", Style: default_style})
				}

				result = append(result, spans...)
				last_line = child.Line

			}
		}
	}

	return result

}

func (renderer *markdownRenderer) inlineElement(element *parser_markdown.Element, default_kind string) []transcriptSpan {

	result := make([]transcriptSpan, 0)
	default_style := markdownStyle(renderer.styles, default_kind)

	switch element.Type {

	case "#text":

		if element.Text != "" {
			result = append(result, transcriptSpan{Text: element.Text, Style: default_style})
		}

	case "b":

		result = append(result, renderer.inlineTextOrChildren(element, "strong")...)

	case "em":

		result = append(result, renderer.inlineTextOrChildren(element, "emphasis")...)

	case "del":

		result = append(result, renderer.inlineTextOrChildren(element, "strike")...)

	case "u":

		result = append(result, renderer.inlineTextOrChildren(element, "text")...)

	case "code":

		result = append(result, transcriptSpan{Text: element.Text, Style: markdownStyle(renderer.styles, "code")})

	case "a":

		label := element.Text

		if label == "" {
			result = append(result, renderer.inlineChildren(element.Children, "link")...)
		} else {
			result = append(result, transcriptSpan{Text: label, Style: markdownStyle(renderer.styles, "link")})
		}

		url := element.GetAttribute("href")

		if url != "" && url != label {
			result = append(result, transcriptSpan{Text: " (" + url + ")", Style: markdownStyle(renderer.styles, "linkUrl")})
		}

	case "img":

		label := element.GetAttribute("alt")

		result = append(result, transcriptSpan{Text: label, Style: markdownStyle(renderer.styles, "imageText")})

		url := element.GetAttribute("src")

		if url != "" {
			result = append(result, transcriptSpan{Text: " (" + url + ")", Style: markdownStyle(renderer.styles, "linkUrl")})
		}

	case "abbr", "span":

		if len(element.Children) > 0 {
			result = append(result, renderer.inlineChildren(element.Children, default_kind)...)
		} else if element.Text != "" {
			result = append(result, transcriptSpan{Text: element.Text, Style: default_style})
		}

	default:

		if len(element.Children) > 0 {
			result = append(result, renderer.inlineChildren(element.Children, default_kind)...)
		} else if element.Text != "" {
			result = append(result, transcriptSpan{Text: element.Text, Style: default_style})
		}

	}

	return result

}

func (renderer *markdownRenderer) inlineTextOrChildren(element *parser_markdown.Element, kind string) []transcriptSpan {

	if len(element.Children) > 0 {
		return renderer.inlineChildren(element.Children, kind)
	} else if element.Text != "" {
		return []transcriptSpan{
			{Text: element.Text, Style: markdownStyle(renderer.styles, kind)},
		}
	} else {
		return make([]transcriptSpan, 0)
	}

}

func (renderer *markdownRenderer) renderDocument(document *parser_markdown.Document) {

	last_line := 0
	previous_heading := false

	for _, element := range document.Body {

		if element == nil || element.Type == "#text" {
			continue
		} else {

			gap := 0

			if last_line == 0 {
				gap = element.Line - 1
			} else {
				gap = element.Line - last_line - 1
			}

			if last_line > 0 && (previous_heading == true || isMarkdownHeading(element) == true) {
				gap = max(gap, 1)
			}

			for index := 0; index < gap; index++ {
				renderer.lines = append(renderer.lines, transcriptLine{})
			}

			renderer.renderBlock(element)
			last_line = markdownEndLine(element)
			previous_heading = isMarkdownHeading(element)

		}
	}

}

func (renderer *markdownRenderer) renderBlock(element *parser_markdown.Element) {

	switch element.Type {

	case "h1", "h2", "h3", "h4", "h5", "h6":

		spans := renderer.inlineChildren(element.Children, "heading")
		renderer.lines = append(renderer.lines, transcriptLine{Spans: compactSpans(spans)})

	case "p":

		spans := renderer.inlineChildren(element.Children, "text")
		renderer.appendWrapped(nil, nil, spans, false)

	case "pre":

		renderer.renderCodeBlock(element)

	case "ul", "ol":

		renderer.renderList(element)

	case "blockquote":

		renderer.renderBlockquote(element)

	case "hr":

		renderer.lines = append(renderer.lines, transcriptLine{
			Spans: []transcriptSpan{
				{Text: strings.Repeat("─", renderer.width), Style: markdownStyle(renderer.styles, "horizontalRule")},
			},
		})

	case "img", "audio":

		spans := renderer.inlineElement(element, "text")
		renderer.appendWrapped(nil, nil, spans, false)

	case "table":

		for _, line := range strings.Split(element.Text, "\n") {
			spans := []transcriptSpan{
				{Text: line, Style: markdownStyle(renderer.styles, "codeBlock")},
			}
			renderer.appendWrapped(nil, nil, spans, false)
		}

	default:

		if len(element.Children) > 0 {
			spans := renderer.inlineChildren(element.Children, "text")
			renderer.appendWrapped(nil, nil, spans, false)
		} else if strings.TrimSpace(element.Text) != "" {
			spans := []transcriptSpan{
				{Text: element.Text, Style: markdownStyle(renderer.styles, "text")},
			}
			renderer.appendWrapped(nil, nil, spans, false)
		}

	}

}

func (renderer *markdownRenderer) renderCodeBlock(element *parser_markdown.Element) {

	lines := strings.Split(element.Text, "\n")

	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		renderer.lines = append(renderer.lines, transcriptLine{})
		return
	}

	for _, line := range lines {

		prefix := []transcriptSpan{
			{Text: "  ", Style: markdownStyle(renderer.styles, "codeBlock")},
		}
		spans := []transcriptSpan{
			{Text: line, Style: markdownStyle(renderer.styles, "codeBlock")},
		}

		renderer.appendWrapped(prefix, prefix, spans, true)

	}

}

func (renderer *markdownRenderer) renderList(element *parser_markdown.Element) {

	ordered := element.Type == "ol"
	index := 0

	for _, child := range element.Children {

		if child == nil || child.Type != "li" {
			continue
		} else {

			index++
			marker := "•"

			if ordered == true {
				marker = strconv.Itoa(index) + "."
			}

			renderer.addList("", marker, child.Children, ordered)

		}
	}

}

func (renderer *markdownRenderer) renderBlockquote(element *parser_markdown.Element) {

	prefix := []transcriptSpan{
		{Text: "│ ", Style: markdownStyle(renderer.styles, "blockQuote")},
	}

	for _, child := range element.Children {

		if child == nil || child.Type == "#text" {
			continue
		} else {

			if child.Type == "p" {
				spans := renderer.inlineChildren(child.Children, "blockQuote")
				renderer.appendWrapped(prefix, prefix, spans, false)
			} else {
				spans := renderer.inlineElement(child, "blockQuote")
				renderer.appendWrapped(prefix, prefix, spans, false)
			}

		}
	}

}
