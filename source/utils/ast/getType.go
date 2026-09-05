package ast

import "bytes"
import "go/ast"
import "go/printer"
import "go/token"
import "strings"

func getType(file *ast.File, fileset *token.FileSet, symbol string, expected_type string) string {

	for _, decl := range file.Decls {

		gen_decl, ok1 := decl.(*ast.GenDecl)

		if ok1 == true && gen_decl.Tok == token.TYPE {

			for _, spec := range gen_decl.Specs {

				type_spec, ok2 := spec.(*ast.TypeSpec)

				if ok2 == true && type_spec.Name != nil && type_spec.Name.Name == symbol {

					if typeName(fileset, type_spec.Type) == expected_type {

						buffer := bytes.Buffer{}
						printer.Fprint(&buffer, fileset, gen_decl)

						return strings.TrimSpace(buffer.String())

					}

				}

			}

		}

	}

	return ""

}
