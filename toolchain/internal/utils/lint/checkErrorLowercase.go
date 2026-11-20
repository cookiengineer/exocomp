package lint

import "go/ast"
import "strings"

func checkErrorLowercase(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	ast.Inspect(file.AST, func(node ast.Node) bool {

		call, ok := node.(*ast.CallExpr)

		if ok == true {

			message, is_literal := errorfLiteral(call)

			if is_literal == true && message != "" && prefixedErrorPattern.MatchString(message) == true {

				message = stripErrorPrefix(message)

				if message != "" {

					first_word := message

					if index := strings.IndexAny(first_word, " \t"); index >= 0 {
						first_word = first_word[:index]
					}

					if strings.ContainsAny(first_word, "-_") == false {

						first := message[0]

						if first >= 'a' && first <= 'z' {

							result = append(
								result,
								file.Diagnostic(
									"errors/lowercase",
									call.Pos(),
									"error %q starts with a lowercase letter",
									message,
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
