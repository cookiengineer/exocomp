package components

func runeWidth(value rune) int {

	if value == '\t' {
		return 4
	}

	if value < 32 || value == 0x7f {
		return 0
	}

	if value >= 0x1100 && value <= 0x115f {
		return 2
	}

	return 1

}
