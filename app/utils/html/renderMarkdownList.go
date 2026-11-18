package html

import parser_markdown "exocomp/parsers/markdown"
import "strings"

func renderMarkdownList(element *parser_markdown.Element) string {

	builder := strings.Builder{}

	builder.WriteString("<" + element.Type + ">")

	for _, child := range element.Children {

		if child == nil || child.Type != "li" {
			continue
		}

		builder.WriteString("<li>")
		builder.WriteString(renderMarkdownInlineChildren(child.Children))

		for _, nested := range child.Children {

			if nested != nil && (nested.Type == "ul" || nested.Type == "ol") {
				builder.WriteString(renderMarkdownList(nested))
			}

		}

		builder.WriteString("</li>")

	}

	builder.WriteString("</" + element.Type + ">")

	return builder.String()

}
