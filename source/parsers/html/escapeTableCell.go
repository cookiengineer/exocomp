package html

import "strings"

func escapeTableCell(value string) string {

	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.ReplaceAll(value, "\n", " ")

	return strings.TrimSpace(value)

}
