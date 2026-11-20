package lint

import "go/token"

// clauseStartLine returns the first line of the comment block attached to the
// clause at pos, or the clause line when no comment is attached.
func clauseStartLine(file *File, pos token.Pos) int {

	line := file.Line(pos)

	for {

		attached := 0

		for _, group := range file.AST.Comments {

			start_line := file.Line(group.Pos())
			end_line := file.Line(group.End())

			if end_line != line-1 || start_line >= line {
				continue
			}

			if start_line-1 < 0 || start_line-1 >= len(file.Lines) {
				continue
			}

			if isCommentLine(file.Lines[start_line-1]) == false {
				continue
			}

			attached = start_line
			break

		}

		if attached == 0 {
			return line
		} else {
			line = attached
		}

	}

}
