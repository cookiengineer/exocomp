package main

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
