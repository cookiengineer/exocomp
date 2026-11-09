package html

func isWhitespace(chr byte) bool {
	return chr == ' ' || chr == '\t' || chr == '\n' || chr == '\r' || chr == '\f' || chr == '\v'
}
