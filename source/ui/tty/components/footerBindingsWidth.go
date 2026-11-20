package components

import utils_tty "exocomp/utils/tty"

func footerBindingsWidth(bindings []footerBinding) int {

	width := 0

	for index, binding := range bindings {

		if index > 0 {
			width += utils_tty.StringWidth(" · ")
		}

		width += utils_tty.StringWidth(binding.Key)
		width += 1
		width += utils_tty.StringWidth(binding.Label)

	}

	return width

}
