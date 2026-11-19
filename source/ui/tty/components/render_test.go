package components

import utils_tty "exocomp/utils/tty"
import "strings"

func renderTestLine(buffer *utils_tty.Buffer, y int, width int) string {

	builder := strings.Builder{}

	for x := 0; x < width; x++ {
		builder.WriteString(buffer.Get(x, y).Text)
	}

	return builder.String()

}
