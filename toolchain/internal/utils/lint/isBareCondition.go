package lint

import "go/ast"
import "go/token"

func isBareCondition(expr ast.Expr) (ast.Node, bool) {

	unwrapped := unwrapParens(expr)

	switch typed := unwrapped.(type) {

	case *ast.Ident:

		if typed.Name == "true" || typed.Name == "false" || typed.Name == "nil" {
			return nil, false
		} else {
			return typed, true
		}

	case *ast.SelectorExpr:
		return typed, true

	case *ast.CallExpr:
		return typed, true

	case *ast.IndexExpr:
		return typed, true

	case *ast.TypeAssertExpr:
		return typed, true

	case *ast.UnaryExpr:

		if typed.Op == token.NOT {
			return nil, false
		} else {
			return typed, true
		}

	default:
		return nil, false

	}

}
