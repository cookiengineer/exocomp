package ast

import "bytes"
import "go/ast"
import "go/parser"
import "go/printer"
import "go/token"
import "strings"

func parseDeclaration(declaration string, declaration_type string) (string, string) {

	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", []byte("package dummy\n"+declaration), 0)

	if err != nil {
		return "", ""
	}

	for _, decl := range file.Decls {

		if declaration_type == "func" {

			func_decl, ok := decl.(*ast.FuncDecl)

			if ok == true && func_decl.Name != nil {

				buffer := bytes.Buffer{}
				printer.Fprint(&buffer, fileset, func_decl)

				name := func_decl.Name.Name

				if func_decl.Recv != nil && len(func_decl.Recv.List) > 0 {
					name = receiverTypeName(fileset, func_decl.Recv.List[0].Type) + "." + name
				}

				return strings.TrimSpace(buffer.String()), name

			}

		}

		gen_decl, ok := decl.(*ast.GenDecl)

		if ok == true && gen_decl.Tok == token.TYPE {

			for _, spec := range gen_decl.Specs {

				type_spec, ok2 := spec.(*ast.TypeSpec)

				if ok2 == true && type_spec.Name != nil {

					if declaration_type != "func" || typeName(fileset, type_spec.Type) == "func" {

						buffer := bytes.Buffer{}
						printer.Fprint(&buffer, fileset, gen_decl)

						return strings.TrimSpace(buffer.String()), type_spec.Name.Name

					}

				}

			}

		}

	}

	return "", ""

}

