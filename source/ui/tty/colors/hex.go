package colors

import utils_tty "exocomp/utils/tty"

func hex(value string) utils_tty.Color {

	color, err := utils_tty.ColorFromHex(value)

	if err == nil {
		return color
	} else {
		panic("colors.hex: " + err.Error())
	}

}
