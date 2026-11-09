package tty

import "unicode"

func RuneWidth(value rune) int {

	if value == 0 {
		return 0
	}

	if value < 32 || (value >= 0x7f && value < 0xa0) {
		return 0
	}

	if unicode.Is(unicode.Mn, value) || unicode.Is(unicode.Me, value) || unicode.Is(unicode.Cf, value) {
		return 0
	}

	if isWideRune(value) {
		return 2
	}

	return 1

}
