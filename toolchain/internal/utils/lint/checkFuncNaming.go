package lint

import "go/ast"
import "strings"

func checkFuncNaming(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, decl := range file.AST.Decls {

			func_decl, ok := decl.(*ast.FuncDecl)

			if ok == false {
				continue
			}

			name := func_decl.Name.Name

			if isExported(name) == true || strings.Contains(name, "_") == false {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"naming/func",
					func_decl.Pos(),
					"unexported function %s uses snake_case; use camelCase",
					name,
				),
			)

		}

		return result

	}

}
