package types

import "unicode"

func skipWhitespace(raw string, index *int) {

	for *index < len(raw) && unicode.IsSpace(rune(raw[*index])) {
		*index++
	}

}
