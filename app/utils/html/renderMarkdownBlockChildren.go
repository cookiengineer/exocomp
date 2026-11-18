package html

import parser_markdown "exocomp/parsers/markdown"
import "strings"

func renderMarkdownBlockChildren(children []*parser_markdown.Element) string {

	builder := strings.Builder{}

	for _, child := range children {

		if child == nil || child.Type == "#text" {
			continue
		}

		rendered := renderMarkdownBlock(child)

		if strings.TrimSpace(rendered) != "" {
			builder.WriteString(rendered)
		}

	}

	return builder.String()

}
