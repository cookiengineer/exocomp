package html

import "strings"

// normalizeText collapses runs of whitespace into a single space, keeping a
// single leading/trailing space when the original text had one, so that
// inline elements stay separated ("Hello <b>world</b>" -> "Hello **world**").
func normalizeText(raw string) string {

	collapsed := strings.Join(strings.Fields(raw), " ")

	if collapsed == "" {
		return ""
	}

	result := collapsed

	if len(raw) > 0 && isWhitespace(raw[0]) {
		result = " " + result
	}

	if len(raw) > 0 && isWhitespace(raw[len(raw)-1]) {
		result = result + " "
	}

	return result

}
