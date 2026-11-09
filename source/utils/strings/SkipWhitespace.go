package strings

import "unicode"

func SkipWhitespace(raw string, index *int) {

	for *index < len(raw) && unicode.IsSpace(rune(raw[*index])) {
		*index++
	}

}
