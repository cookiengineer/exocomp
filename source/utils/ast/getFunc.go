package ast

import "bytes"
import "go/ast"
import "go/printer"
import "go/token"
import "strings"

func getFunc(file *ast.File, fileset *token.FileSet, symbol string) *Symbol {

	for _, decl := range file.Decls {

		func_decl, ok := decl.(*ast.FuncDecl)

		if ok == true && func_decl.Name != nil {

			if func_decl.Recv != nil && len(func_decl.Recv.List) > 0 {

				receiver_type := receiverTypeName(fileset, func_decl.Recv.List[0].Type)

				if receiver_type+"."+func_decl.Name.Name == symbol {

					buffer := bytes.Buffer{}
					printer.Fprint(&buffer, fileset, func_decl)

					return &Symbol{
						Name: receiver_type+"."+func_decl.Name.Name,
						Type: "func",
						Body: strings.TrimSpace(buffer.String()),
					}

				}

			} else if func_decl.Name.Name == symbol {

				buffer := bytes.Buffer{}
				printer.Fprint(&buffer, fileset, func_decl)

				return &Symbol{
					Name: func_decl.Name.Name,
					Type: "func",
					Body: strings.TrimSpace(buffer.String()),
				}

			}

		}

	}

	return nil

}
