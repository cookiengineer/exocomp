package lint

import "go/ast"
import "go/token"

func checkLocalCase(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		report := func(name *ast.Ident) {

			if name == nil || name.Name == "_" {
				return
			}

			if isCamelLocal(name.Name) == false {
				return
			}

			result = append(
				result,
				file.Diagnostic(
					"style/local-case",
					name.Pos(),
					"local variable %s uses camelCase; use snake_case",
					name.Name,
				),
			)

		}

		report_list := func(fields *ast.FieldList) {

			if fields == nil {
				return
			}

			for _, field := range fields.List {
				for _, name := range field.Names {
					report(name)
				}
			}

		}

		for _, decl := range file.AST.Decls {

			func_decl, ok := decl.(*ast.FuncDecl)

			if ok == false || func_decl.Body == nil {
				continue
			}

			report_list(func_decl.Type.Params)

			ast.Inspect(func_decl.Body, func(node ast.Node) bool {

				switch typed := node.(type) {

				case *ast.FuncLit:
					report_list(typed.Type.Params)

				case *ast.AssignStmt:

					if typed.Tok == token.DEFINE {

						for _, left := range typed.Lhs {

							ident, is_ident := left.(*ast.Ident)

							if is_ident == true {
								report(ident)
							}

						}

					}

				case *ast.RangeStmt:

					if typed.Tok == token.DEFINE {

						if ident, is_ident := typed.Key.(*ast.Ident); is_ident == true {
							report(ident)
						}

						if ident, is_ident := typed.Value.(*ast.Ident); is_ident == true {
							report(ident)
						}

					}

				case *ast.ValueSpec:

					for _, name := range typed.Names {
						report(name)
					}

				}

				return true

			})

		}

		return result

	}

}
