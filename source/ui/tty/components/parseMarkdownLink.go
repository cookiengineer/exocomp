package components

import "strings"

func parseMarkdownLink(text string, start int, image bool) (string, string, int, bool) {

	offset := start + 1

	if image == true {

		if strings.HasPrefix(text[start:], "![") == false {
			return "", "", start, false
		}

		offset = start + 2

	}

	close_bracket := strings.Index(text[offset:], "](")

	if close_bracket < 0 {
		return "", "", start, false
	}

	label := text[offset : offset+close_bracket]
	url_start := offset + close_bracket + 2
	close_parenthesis := strings.Index(text[url_start:], ")")

	if close_parenthesis < 0 {
		return "", "", start, false
	}

	url := text[url_start : url_start+close_parenthesis]

	return label, url, url_start + close_parenthesis + 1, true

}
