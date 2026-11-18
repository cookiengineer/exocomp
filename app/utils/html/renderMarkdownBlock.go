package html

import parser_markdown "exocomp/parsers/markdown"
import std_html "html"
import "strings"

func renderMarkdownBlock(element *parser_markdown.Element) string {

	switch element.Type {

	case "h1", "h2", "h3", "h4", "h5", "h6":

		return "<" + element.Type + ">" + renderMarkdownInlineChildren(element.Children) + "</" + element.Type + ">"

	case "p":

		return "<p>" + renderMarkdownInlineChildren(element.Children) + "</p>"

	case "pre":

		class := strings.TrimSpace(element.GetAttribute("class"))

		if class != "" {
			return "<pre class=\"" + std_html.EscapeString(class) + "\">" + std_html.EscapeString(element.Text) + "</pre>"
		} else {
			return "<pre>" + std_html.EscapeString(element.Text) + "</pre>"
		}

	case "ul", "ol":

		return renderMarkdownList(element)

	case "blockquote":

		return "<blockquote>" + renderMarkdownBlockChildren(element.Children) + "</blockquote>"

	case "hr":

		return "<hr>"

	case "img", "audio":

		return renderMarkdownInline(element)

	case "table":

		return "<pre>" + std_html.EscapeString(element.Text) + "</pre>"

	default:

		if len(element.Children) > 0 {
			return "<p>" + renderMarkdownInlineChildren(element.Children) + "</p>"
		} else if strings.TrimSpace(element.Text) != "" {
			return "<p>" + std_html.EscapeString(element.Text) + "</p>"
		} else {
			return ""
		}

	}

}
