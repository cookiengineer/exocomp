package lint

import "go/token"
import "strings"

// hasBlankLineAbove reports whether the clause header at pos is preceded by a
// blank line. Comments directly above the clause stay attached to it.
func hasBlankLineAbove(file *File, pos token.Pos) bool {

	start_line := clauseStartLine(file, pos)

	if start_line <= 1 {
		return false
	} else {

		index := start_line - 2

		if index < 0 || index >= len(file.Lines) {
			return false
		} else {
			return strings.TrimSpace(file.Lines[index]) == ""
		}

	}

}
