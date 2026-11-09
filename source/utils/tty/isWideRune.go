package tty

func isWideRune(value rune) bool {

	return (value >= 0x1100 && value <= 0x115f) ||
		value == 0x2329 ||
		value == 0x232a ||
		(value >= 0x2e80 && value <= 0xa4cf && value != 0x303f) ||
		(value >= 0xac00 && value <= 0xd7a3) ||
		(value >= 0xf900 && value <= 0xfaff) ||
		(value >= 0xfe10 && value <= 0xfe19) ||
		(value >= 0xfe30 && value <= 0xfe6f) ||
		(value >= 0xff00 && value <= 0xff60) ||
		(value >= 0xffe0 && value <= 0xffe6) ||
		(value >= 0x1f300 && value <= 0x1f64f) ||
		(value >= 0x1f900 && value <= 0x1f9ff) ||
		(value >= 0x20000 && value <= 0x3fffd)

}
