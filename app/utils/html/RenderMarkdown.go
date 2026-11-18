package html

import parser_markdown "exocomp/parsers/markdown"
import std_html "html"
import "strings"

// RenderMarkdown parses markdown and renders safe HTML for a chat message.
func RenderMarkdown(raw string) string {

	document := parser_markdown.ParseFragment(raw)
	rendered := RenderDocument(document)

	if rendered == "" && strings.TrimSpace(raw) != "" {
		return "<p>" + std_html.EscapeString(strings.TrimSpace(raw)) + "</p>"
	}

	return rendered

}
