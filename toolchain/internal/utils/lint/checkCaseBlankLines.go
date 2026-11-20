package lint

import "go/ast"

// checkCaseBlankLines requires a blank line above every case and default clause.
func checkCaseBlankLines(file *File, block *ast.BlockStmt) []Diagnostic {

	result := make([]Diagnostic, 0)

	for _, stmt := range block.List {

		case_clause, ok := stmt.(*ast.CaseClause)

		if ok == false {
			continue
		}

		if hasBlankLineAbove(file, case_clause.Pos()) == true {
			continue
		}

		result = append(
			result,
			file.Diagnostic(
				"style/blank-block",
				case_clause.Pos(),
				"case clause needs a blank line above it",
			),
		)

	}

	return result

}
