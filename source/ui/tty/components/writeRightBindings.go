package components

import "exocomp/ui/tty/styles"
import utils_tty "exocomp/utils/tty"

func writeRightBindings(buffer *utils_tty.Buffer, right int, y int, minimum int, bindings []footerBinding) {

	width := footerBindingsWidth(bindings)
	available := max(0, right-minimum)
	x := max(minimum, right-width)

	if width > available {
		x = minimum
	}

	for index, binding := range bindings {

		if index > 0 {
			x += writeFooterSpan(buffer, x, y, " · ", styles.FooterHint, right)
		}

		x += writeFooterSpan(buffer, x, y, binding.Key, styles.FooterKey, right)
		x += writeFooterSpan(buffer, x, y, " ", styles.FooterHint, right)
		x += writeFooterSpan(buffer, x, y, binding.Label, styles.FooterHint, right)

		if x >= right {
			break
		}

	}

}
