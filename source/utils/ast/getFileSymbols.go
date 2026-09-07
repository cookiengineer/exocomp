package ast

import "bytes"
import "go/ast"
import "go/printer"
import "go/token"
import "strings"

func getFileSymbols(file *ast.File, fileset *token.FileSet, with_private bool) map[string]*Symbol {

	result := make(map[string]*Symbol, 0)

	for _, decl := range file.Decls {

		func_decl, ok1 := decl.(*ast.FuncDecl)

		if ok1 == true && func_decl.Name != nil {

			if func_decl.Recv != nil && len(func_decl.Recv.List) > 0 {

				receiver_type := receiverTypeName(fileset, func_decl.Recv.List[0].Type)

				buffer := bytes.Buffer{}
				printer.Fprint(&buffer, fileset, func_decl)

				declaration := strings.TrimSpace(buffer.String())

				if strings.Contains(declaration, "{\n") {
					declaration = strings.TrimSpace(declaration[0:strings.Index(declaration, "{\n")])
				}

				if isExportedSymbol(receiver_type) || with_private == true {

					result[receiver_type+"."+func_decl.Name.Name] = &Symbol{
						Name: receiver_type+"."+func_decl.Name.Name,
						Type: "func",
						Body: declaration,
					}

				}

			} else {

				buffer := bytes.Buffer{}
				printer.Fprint(&buffer, fileset, func_decl)

				declaration := strings.TrimSpace(buffer.String())

				if strings.Contains(declaration, "{\n") {
					declaration = strings.TrimSpace(declaration[0:strings.Index(declaration, "{\n")])
				}

				if isExportedSymbol(func_decl.Name.Name) || with_private == true {

					result[func_decl.Name.Name] = &Symbol{
						Name: func_decl.Name.Name,
						Type: "func",
						Body: declaration,
					}

				}

			}

		} else {

			gen_decl, ok2 := decl.(*ast.GenDecl)

			if ok2 == true && gen_decl.Tok == token.TYPE {

				for _, spec := range gen_decl.Specs {

					type_spec, ok3 := spec.(*ast.TypeSpec)

					if ok3 == true && type_spec.Name != nil {

						buffer := bytes.Buffer{}
						printer.Fprint(&buffer, fileset, gen_decl)

						declaration := strings.TrimSpace(buffer.String())

						if strings.Contains(declaration, "{\n") {
							declaration = strings.TrimSpace(declaration[0:strings.Index(declaration, "{\n")])
						}

						if isExportedSymbol(type_spec.Name.Name) || with_private == true {

							result[type_spec.Name.Name] = &Symbol{
								Name: type_spec.Name.Name,
								Type: "type",
								Body: declaration,
							}

						}

					}

				}

			}

		}

	}

	return result

}
