package lint

import "go/ast"

func checkBooleanNaming(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		report := func(name *ast.Ident) {

			if name == nil || name.Name == "_" {
				return
			}

			if isPredicateName(name.Name) == false {
				return
			}

			result = append(
				result,
				file.Diagnostic(
					"naming/boolean",
					name.Pos(),
					"boolean %s does not read as a predicate; use is_/has_/can_/should_",
					name.Name,
				),
			)

		}

		ast.Inspect(file.AST, func(node ast.Node) bool {

			switch typed := node.(type) {

			case *ast.FuncDecl:

				if typed.Type.Params != nil {

					for _, field := range typed.Type.Params.List {

						if isBoolType(field.Type) == false {
							continue
						}

						for _, name := range field.Names {
							report(name)
						}

					}

				}

			case *ast.ValueSpec:

				if isBoolType(typed.Type) == true {

					for _, name := range typed.Names {
						report(name)
					}

				}

			}

			return true

		})

		return result

	}

}
