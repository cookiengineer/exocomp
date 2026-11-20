package lint

import "go/ast"

func checkBareCondition(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	report := func(expr ast.Expr) {

		node, ok := isBareCondition(expr)

		if ok == false {
			return
		}

		result = append(
			result,
			file.Diagnostic(
				"cond/bare",
				node.Pos(),
				"condition %q is not compared explicitly; use == true or == false",
				file.Text(node),
			),
		)

	}

	ast.Inspect(file.AST, func(node ast.Node) bool {

		switch typed := node.(type) {

		case *ast.IfStmt:
			report(typed.Cond)

		case *ast.ForStmt:
			report(typed.Cond)

		}

		return true

	})

	return result

}
