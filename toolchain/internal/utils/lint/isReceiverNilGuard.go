package lint

import "go/ast"
import "go/token"

// isReceiverNilGuard reports whether stmt is a `receiver == nil` guard. The
// codestyle permits early returns for nil receivers.
func isReceiverNilGuard(stmt ast.Stmt, receiver string) bool {

	if receiver == "" {
		return false
	} else {

		if_stmt, ok := stmt.(*ast.IfStmt)

		if ok == false {
			return false
		} else {

			binary, is_binary := unwrapParens(if_stmt.Cond).(*ast.BinaryExpr)

			if is_binary == false || (binary.Op != token.EQL && binary.Op != token.NEQ) {
				return false
			} else {

				left, left_ok := unwrapParens(binary.X).(*ast.Ident)
				right, right_ok := unwrapParens(binary.Y).(*ast.Ident)

				if left_ok == false || right_ok == false {
					return false
				} else if left.Name == receiver && right.Name == "nil" {
					return true
				} else {
					return right.Name == receiver && left.Name == "nil"
				}

			}

		}

	}

}
