package lint

import "go/ast"
import "go/token"

// guardKind classifies a block that always leaves the enclosing flow.
// It returns "" when the block is not a terminal guard.
func guardKind(block *ast.BlockStmt) string {

	if len(block.List) == 0 {
		return ""
	} else {

		kind := ""

		for _, stmt := range block.List {

			switch typed := stmt.(type) {

			case *ast.ReturnStmt:

				if len(typed.Results) > 0 {

					if kind == "" || kind == "bare" {
						kind = "value"
					}

				} else if kind == "" {
					kind = "bare"
				}

			case *ast.BranchStmt:

				if typed.Tok == token.FALLTHROUGH {
					return ""
				} else if kind == "" {
					kind = "loop"
				}

			case *ast.ExprStmt:

				call, ok := typed.X.(*ast.CallExpr)

				if ok == false {
					return ""
				} else if isFatalCall(call) == false {
					return ""
				} else if kind == "" {
					kind = "fatal"
				}

			default:
				return ""

			}

		}

		return kind

	}

}
