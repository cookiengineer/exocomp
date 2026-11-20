package lint

import "go/ast"
import "go/token"

func checkNot(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	report := func(expr ast.Expr) {

		if expr == nil {
			return
		}

		ast.Inspect(expr, func(node ast.Node) bool {

			unary, ok := node.(*ast.UnaryExpr)

			if ok == true && unary.Op == token.NOT {
				result = append(
					result,
					file.Diagnostic(
						"cond/not",
						unary.Pos(),
						"do not use !%s; compare explicitly with == false",
						file.Text(unary.X),
					),
				)
			}

			return true

		})

	}

	ast.Inspect(file.AST, func(node ast.Node) bool {

		switch typed := node.(type) {

		case *ast.IfStmt:
			report(typed.Cond)

		case *ast.ForStmt:
			report(typed.Cond)

		case *ast.AssignStmt:

			for _, value := range typed.Rhs {
				report(value)
			}

		case *ast.ValueSpec:

			for _, value := range typed.Values {
				report(value)
			}

		}

		return true

	})

	return result

}
