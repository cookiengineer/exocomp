package lint

import "go/ast"
import "go/token"
import "strings"

// errorfExemptions returns the positions of fmt.Errorf calls whose message is
// prefixed by the caller: unexported helpers inside a tool package. Test
// support binaries under tools/testutil never reach the model at all.
func errorfExemptions(file *File) map[token.Pos]bool {

	exempt := make(map[token.Pos]bool)

	if strings.HasPrefix(file.Rel, "tools/") == false {
		return exempt
	} else if strings.HasPrefix(file.Rel, "tools/testutil/") == true {

		ast.Inspect(file.AST, func(node ast.Node) bool {

			call, ok := node.(*ast.CallExpr)

			if ok == true {
				exempt[call.Pos()] = true
			}

			return true

		})

		return exempt

	} else {

		for _, decl := range file.AST.Decls {

			func_decl, ok := decl.(*ast.FuncDecl)

			if ok == false || func_decl.Body == nil {
				continue
			}

			if isExported(func_decl.Name.Name) == true {
				continue
			}

			ast.Inspect(func_decl.Body, func(node ast.Node) bool {

				call, ok := node.(*ast.CallExpr)

				if ok == true {
					exempt[call.Pos()] = true
				}

				return true

			})

		}

		return exempt

	}

}
