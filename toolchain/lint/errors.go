package main

import "go/ast"
import "go/token"
import "regexp"
import "strconv"
import "strings"

var prefixedErrorPattern = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9_]*\.)+(%s|%v|[A-Z][A-Za-z0-9]*):`)

var parameterPrefixPattern = regexp.MustCompile(`^%(s|v)([.:]|$)`)

func errorfLiteral(call *ast.CallExpr) (string, bool) {

	selector, ok := call.Fun.(*ast.SelectorExpr)

	if ok == false || selector.Sel.Name != "Errorf" {
		return "", false
	} else {

		ident, ok := selector.X.(*ast.Ident)

		if ok == false || ident.Name != "fmt" {
			return "", false
		} else if len(call.Args) == 0 {
			return "", false
		} else {

			literal, ok := call.Args[0].(*ast.BasicLit)

			if ok == false || literal.Kind != token.STRING {
				return "", false
			} else {

				value, err := strconv.Unquote(literal.Value)

				if err == nil {
					return value, true
				} else {
					return "", false
				}

			}

		}

	}

}

func isPrefixedError(message string) bool {

	if prefixedErrorPattern.MatchString(message) == true {
		return true
	} else {
		return parameterPrefixPattern.MatchString(message)
	}

}

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

// stripErrorPrefix removes a leading component prefix.
func stripErrorPrefix(message string) string {

	index := strings.Index(message, ": ")

	if index < 0 {
		return message
	} else {

		head := message[:index]

		if prefixedErrorPattern.MatchString(head+":") == true {
			return message[index+2:]
		} else {
			return message
		}

	}

}
