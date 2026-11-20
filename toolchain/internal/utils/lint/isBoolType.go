package lint

import "go/ast"

func isBoolType(expr ast.Expr) bool {

	if expr == nil {
		return false
	} else {

		ident, ok := expr.(*ast.Ident)

		return ok == true && ident.Name == "bool"

	}

}
