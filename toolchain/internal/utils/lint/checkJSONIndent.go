package lint

import "go/ast"
import "go/token"
import "strconv"

func checkJSONIndent(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	ast.Inspect(file.AST, func(node ast.Node) bool {

		call, ok := node.(*ast.CallExpr)

		if ok == true {

			selector, is_selector := call.Fun.(*ast.SelectorExpr)

			if is_selector == true && selector.Sel.Name == "MarshalIndent" {

				ident, is_ident := selector.X.(*ast.Ident)

				if is_ident == true && ident.Name == "json" && len(call.Args) >= 3 {

					literal, is_literal := call.Args[2].(*ast.BasicLit)

					if is_literal == true && literal.Kind == token.STRING {

						value, err := strconv.Unquote(literal.Value)

						if err == nil && value != "\t" {

							result = append(
								result,
								file.Diagnostic(
									"style/json-indent",
									call.Pos(),
									"JSON indentation must be a tab, got %q",
									value,
								),
							)

						}

					}

				}

			}

		}

		return true

	})

	return result

}
