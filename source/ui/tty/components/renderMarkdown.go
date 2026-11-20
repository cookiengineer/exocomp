package components

import parser_markdown "exocomp/parsers/markdown"
import "strings"

func renderMarkdown(text string, width int) []transcriptLine {

	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	document := parser_markdown.ParseFragment(normalized)

	if document == nil || (len(document.Errors) > 0 && len(document.Body) == 0) {
		return renderMarkdownPlain(normalized, width)
	} else {

		renderer := newMarkdownRenderer(width)
		renderer.renderDocument(document)

		return renderer.lines

	}
}
