package tty

func decodeModifier(value int) (bool, bool, bool) {

	if value <= 0 {
		return false, false, false
	}

	bits := value - 1

	return bits&1 != 0, bits&2 != 0, bits&4 != 0

}
