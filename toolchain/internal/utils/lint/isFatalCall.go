package lint

import "go/ast"
import "strings"

func isFatalCall(call *ast.CallExpr) bool {

	switch typed := call.Fun.(type) {

	case *ast.Ident:
		return typed.Name == "panic"

	case *ast.SelectorExpr:

		name := typed.Sel.Name

		if name == "Exit" {
			return true
		} else {
			return strings.HasPrefix(name, "Fatal") == true || strings.HasPrefix(name, "Skip") == true
		}

	default:
		return false

	}

}
