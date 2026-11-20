package lint

import "go/ast"

func checkAny(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	ast.Inspect(file.AST, func(node ast.Node) bool {

		interface_type, ok := node.(*ast.InterfaceType)

		if ok == true && interface_type.Methods != nil && len(interface_type.Methods.List) == 0 {

			result = append(
				result,
				file.Diagnostic(
					"style/any",
					interface_type.Pos(),
					"use any instead of interface{}",
				),
			)

		}

		return true

	})

	return result

}
