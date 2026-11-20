package lint

import "go/ast"

// collectSwitchBodies returns the body blocks of every switch and type switch.
func collectSwitchBodies(file *File) map[*ast.BlockStmt]bool {

	result := make(map[*ast.BlockStmt]bool)

	ast.Inspect(file.AST, func(node ast.Node) bool {

		switch typed := node.(type) {

		case *ast.SwitchStmt:
			result[typed.Body] = true

		case *ast.TypeSwitchStmt:
			result[typed.Body] = true

		}

		return true

	})

	return result

}
