package html

import parser_markdown "exocomp/parsers/markdown"
import std_html "html"

func renderMarkdownInline(element *parser_markdown.Element) string {

	switch element.Type {

	case "#text":

		return std_html.EscapeString(element.Text)

	case "b":

		return "<strong>" + std_html.EscapeString(element.Text) + "</strong>"

	case "em":

		return "<em>" + std_html.EscapeString(element.Text) + "</em>"

	case "del":

		return "<del>" + std_html.EscapeString(element.Text) + "</del>"

	case "u":

		return "<u>" + std_html.EscapeString(element.Text) + "</u>"

	case "code":

		return "<code>" + std_html.EscapeString(element.Text) + "</code>"

	case "a":

		href := sanitizeMarkdownURL(element.GetAttribute("href"))
		label := std_html.EscapeString(element.Text)

		if label == "" {
			label = renderMarkdownInlineChildren(element.Children)
		}

		return "<a href=\"" + std_html.EscapeString(href) + "\" target=\"_blank\" rel=\"noopener noreferrer\">" + label + "</a>"

	case "abbr":

		title := std_html.EscapeString(element.GetAttribute("title"))

		if title != "" {
			return "<abbr title=\"" + title + "\">" + std_html.EscapeString(element.Text) + "</abbr>"
		} else {
			return std_html.EscapeString(element.Text)
		}

	case "img":

		src := sanitizeMarkdownURL(element.GetAttribute("src"))
		alt := std_html.EscapeString(element.GetAttribute("alt"))

		return "<img src=\"" + std_html.EscapeString(src) + "\" alt=\"" + alt + "\">"

	case "audio":

		src := sanitizeMarkdownURL(element.GetAttribute("src"))
		title := std_html.EscapeString(element.GetAttribute("title"))

		return "<audio controls src=\"" + std_html.EscapeString(src) + "\" title=\"" + title + "\"></audio>"

	case "span":

		return renderMarkdownInlineChildren(element.Children)

	default:

		if len(element.Children) > 0 {
			return renderMarkdownInlineChildren(element.Children)
		} else {
			return std_html.EscapeString(element.Text)
		}

	}

}
