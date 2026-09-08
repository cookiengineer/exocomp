package ast

import "go/ast"
import "go/token"

func writeType(source []byte, fileset *token.FileSet, file *ast.File, symbol string, expected_type string, replacement string) ([]byte, bool) {

	for _, decl := range file.Decls {

		gen_decl, ok1 := decl.(*ast.GenDecl)

		if ok1 == true && gen_decl.Tok == token.TYPE {

			for _, spec := range gen_decl.Specs {

				type_spec, ok2 := spec.(*ast.TypeSpec)

				if ok2 == true && type_spec.Name != nil && type_spec.Name.Name == symbol {

					if expected_type == "" || typeName(fileset, type_spec.Type) == expected_type {
						return spliceNode(source, fileset, gen_decl, replacement), true
					}

				}

			}

		}

	}

	return source, false

}

