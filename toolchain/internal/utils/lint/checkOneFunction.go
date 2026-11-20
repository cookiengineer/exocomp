package lint

import "go/ast"

func checkOneFunction(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true || file.AST.Name.Name == "main" {
		return nil
	} else if oneFunctionExempt[file.Rel] == true {
		return nil
	} else {

		candidates := make([]*ast.FuncDecl, 0)

		for _, func_decl := range freeFunctions(file) {

			if ownedTypeName(func_decl, file.Pkg.Types) == "" {
				candidates = append(candidates, func_decl)
			}

		}

		if len(candidates) <= 1 {
			return nil
		} else {

			base := identifierBaseName(file)
			matching := 0

			for _, func_decl := range candidates {

				if func_decl.Name.Name == base {
					matching++
				}

			}

			result := make([]Diagnostic, 0)

			for _, func_decl := range candidates {

				if matching == 1 && func_decl.Name.Name == base {
					continue
				}

				result = append(
					result,
					file.Diagnostic(
						"layout/one-function",
						func_decl.Pos(),
						"file declares %d free functions; move %s into %s.go",
						len(candidates),
						func_decl.Name.Name,
						func_decl.Name.Name,
					),
				)

			}

			return result

		}

	}

}
