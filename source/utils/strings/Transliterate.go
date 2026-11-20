package strings

// Transliterate returns the ASCII replacement for a single rune. It returns an
// empty string when the rune has no known replacement.
func Transliterate(character rune) string {
	return lookup(uint16(character))
}
