package lint

import "go/ast"

func receiverTypeName(expr ast.Expr) string {

	switch typed := expr.(type) {

	case *ast.StarExpr:
		return receiverTypeName(typed.X)

	case *ast.Ident:
		return typed.Name

	case *ast.IndexExpr:
		return receiverTypeName(typed.X)

	case *ast.IndexListExpr:
		return receiverTypeName(typed.X)

	default:
		return ""

	}

}
