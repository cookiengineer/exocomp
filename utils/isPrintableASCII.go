package utils

func isPrintableASCII(chr byte) bool {
	return chr >= 0x20 && chr <= 0x7E
}
