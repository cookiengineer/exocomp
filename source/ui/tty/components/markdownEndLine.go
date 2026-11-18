package components

import parser_markdown "exocomp/parsers/markdown"
import "strings"

func markdownEndLine(element *parser_markdown.Element) int {

	if element == nil {
		return 0
	}

	end := element.Line

	if element.Type == "pre" {
		return element.Line + strings.Count(element.Text, "\n") + 2
	}

	for _, child := range element.Children {

		if child == nil {
			continue
		}

		child_end := markdownEndLine(child)

		if child_end > end {
			end = child_end
		}

	}

	return end

}
