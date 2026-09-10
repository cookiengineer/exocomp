package ast

import "bytes"
import "go/ast"
import "go/parser"
import "go/printer"
import "go/token"
import "strings"

func SearchSymbols(source []byte, query string) []*Symbol {

	result := make([]*Symbol, 0)
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", source, 0)

	if err == nil {

		for _, decl := range file.Decls {

			func_decl, ok1 := decl.(*ast.FuncDecl)

			if ok1 == true && func_decl.Name != nil {

				name := func_decl.Name.Name

				if func_decl.Recv != nil && len(func_decl.Recv.List) > 0 {
					name = receiverTypeName(fileset, func_decl.Recv.List[0].Type) + "." + func_decl.Name.Name
				}

				buffer := bytes.Buffer{}
				printer.Fprint(&buffer, fileset, func_decl)

				if strings.Contains(strings.ToLower(buffer.String()), strings.ToLower(query)) == true {
					result = append(result, &Symbol{
						Name: name,
						Type: "func",
						Body: strings.TrimSpace(buffer.String()),
					})
				}

			} else {

				gen_decl, ok2 := decl.(*ast.GenDecl)

				if ok2 == true && gen_decl.Tok == token.TYPE {

					for _, spec := range gen_decl.Specs {

						type_spec, ok3 := spec.(*ast.TypeSpec)

						if ok3 == true && type_spec.Name != nil {

							buffer := bytes.Buffer{}
							printer.Fprint(&buffer, fileset, gen_decl)

							if strings.Contains(strings.ToLower(buffer.String()), strings.ToLower(query)) == true {
								result = append(result, &Symbol{
									Name: type_spec.Name.Name,
									Type: "type",
									Body: strings.TrimSpace(buffer.String()),
								})
							}

						}

					}

				}

			}

		}

	}

	return result

}
