package lint

import "go/ast"

func isValueGuard(stmt ast.Stmt) bool {

	if_stmt, ok := stmt.(*ast.IfStmt)

	if ok == false || if_stmt.Else != nil {
		return false
	} else {
		return guardKind(if_stmt.Body) == "value"
	}

}
