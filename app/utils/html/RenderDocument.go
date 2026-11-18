package html

import parser_markdown "exocomp/parsers/markdown"
import "strings"

func RenderDocument(document *parser_markdown.Document) string {

	if document == nil {
		return ""
	}

	builder := strings.Builder{}

	for _, element := range document.Body {

		if element == nil || element.Type == "#text" {
			continue
		}

		rendered := renderMarkdownBlock(element)

		if strings.TrimSpace(rendered) != "" {
			builder.WriteString(rendered)
		}

	}

	return builder.String()

}
