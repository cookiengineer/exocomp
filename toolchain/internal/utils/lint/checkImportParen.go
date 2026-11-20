package lint

import "go/ast"
import "go/token"

func checkImportParen(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	for _, decl := range file.AST.Decls {

		gen_decl, ok := decl.(*ast.GenDecl)

		if ok == false || gen_decl.Tok != token.IMPORT {
			continue
		}

		if gen_decl.Lparen.IsValid() == true {
			result = append(
				result,
				file.Diagnostic(
					"import/paren",
					gen_decl.Pos(),
					"imports must be one per line, without a parenthesized block",
				),
			)
		}

	}

	return result

}
