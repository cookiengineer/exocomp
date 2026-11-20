package lint

import "go/ast"

func checkFunctionFile(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true || file.AST.Name.Name == "main" {
		return nil
	} else {

		candidates := make([]*ast.FuncDecl, 0)

		for _, func_decl := range freeFunctions(file) {

			if ownedTypeName(func_decl, file.Pkg.Types) == "" {
				candidates = append(candidates, func_decl)
			}

		}

		if len(candidates) != 1 {
			return nil
		} else {

			base := identifierBaseName(file)

			if candidates[0].Name.Name == base {
				return nil
			} else {
				return []Diagnostic{
					file.Diagnostic(
						"layout/function-file",
						candidates[0].Pos(),
						"%s must live in %s.go",
						candidates[0].Name.Name,
						candidates[0].Name.Name,
					),
				}
			}

		}

	}

}
