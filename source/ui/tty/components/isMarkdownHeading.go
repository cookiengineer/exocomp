package components

import parser_markdown "exocomp/parsers/markdown"

func isMarkdownHeading(element *parser_markdown.Element) bool {

	if element != nil {

		switch element.Type {

		case "h1", "h2", "h3", "h4", "h5", "h6":
			return true

		default:
			return false

		}

	} else {
		return false
	}
}
