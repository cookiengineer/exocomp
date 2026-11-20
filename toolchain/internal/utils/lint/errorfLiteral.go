package lint

import "go/ast"
import "go/token"
import "strconv"

func errorfLiteral(call *ast.CallExpr) (string, bool) {

	selector, ok := call.Fun.(*ast.SelectorExpr)

	if ok == false || selector.Sel.Name != "Errorf" {
		return "", false
	} else {

		ident, ok := selector.X.(*ast.Ident)

		if ok == false || ident.Name != "fmt" {
			return "", false
		} else if len(call.Args) == 0 {
			return "", false
		} else {

			literal, ok := call.Args[0].(*ast.BasicLit)

			if ok == false || literal.Kind != token.STRING {
				return "", false
			} else {

				value, err := strconv.Unquote(literal.Value)

				if err == nil {
					return value, true
				} else {
					return "", false
				}

			}

		}

	}

}
