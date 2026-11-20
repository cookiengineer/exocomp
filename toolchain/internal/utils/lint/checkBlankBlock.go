package lint

import "go/ast"
import "strings"

func checkBlankBlock(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)
	switch_bodies := collectSwitchBodies(file)

	ast.Inspect(file.AST, func(node ast.Node) bool {

		block, ok := node.(*ast.BlockStmt)

		if ok == true {

			is_switch_body := switch_bodies[block]

			if is_switch_body == true {
				result = append(result, checkCaseBlankLines(file, block)...)
			}

			if len(block.List) > 1 {

				open_line := file.Line(block.Lbrace)
				close_line := file.Line(block.Rbrace)

				if open_line != close_line {

					after := open_line + 1

					if is_switch_body == false && after <= close_line && after-1 < len(file.Lines) {

						line := file.Lines[after-1]

						if strings.TrimSpace(line) != "" && isCommentLine(line) == false {
							result = append(
								result,
								file.Diagnostic(
									"style/blank-block",
									block.Lbrace,
									"multi-statement block needs a blank line after the opening brace",
								),
							)
						}

					}

					before := close_line - 1

					if before > open_line && before-1 < len(file.Lines) && before-1 >= 0 {

						line := file.Lines[before-1]

						if strings.TrimSpace(line) != "" && isCommentLine(line) == false {
							result = append(
								result,
								file.Diagnostic(
									"style/blank-block",
									block.Rbrace,
									"multi-statement block needs a blank line before the closing brace",
								),
							)
						}

					}

				}

			}

		}

		return true

	})

	return result

}
