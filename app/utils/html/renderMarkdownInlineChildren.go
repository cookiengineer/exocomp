package html

import parser_markdown "exocomp/parsers/markdown"
import "strings"

func renderMarkdownInlineChildren(children []*parser_markdown.Element) string {

	builder := strings.Builder{}
	last_line := 0

	for _, child := range children {

		if child == nil {
			continue
		}

		if child.Type == "ul" || child.Type == "ol" || child.Type == "pre" {
			continue
		}

		rendered := renderMarkdownInline(child)

		if rendered == "" {
			continue
		}

		if last_line > 0 && child.Line > last_line {
			builder.WriteString("<br>")
		} else if builder.Len() > 0 && needsMarkdownSpace(builder.String(), rendered) == true {
			builder.WriteString(" ")
		}

		builder.WriteString(rendered)
		last_line = child.Line

	}

	return builder.String()

}
