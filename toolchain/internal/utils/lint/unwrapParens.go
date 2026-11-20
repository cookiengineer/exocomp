package lint

import "go/ast"

func unwrapParens(expr ast.Expr) ast.Expr {

	for {

		typed, ok := expr.(*ast.ParenExpr)

		if ok == true {
			expr = typed.X
		} else {
			return expr
		}

	}

}
