package lint

import "go/ast"

func checkErrorPrefix(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)
	exempt := errorfExemptions(file)

	ast.Inspect(file.AST, func(node ast.Node) bool {

		call, ok := node.(*ast.CallExpr)

		if ok == true && exempt[call.Pos()] == false {

			message, is_literal := errorfLiteral(call)

			if is_literal == true && isPrefixedError(message) == false {

				result = append(
					result,
					file.Diagnostic(
						"errors/prefix",
						call.Pos(),
						"error %q is missing a component prefix such as files.Read:",
						message,
					),
				)

			}

		}

		return true

	})

	return result

}
