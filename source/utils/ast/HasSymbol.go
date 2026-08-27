package ast

import "bytes"
import "go/ast"
import "go/parser"
import "go/printer"
import "go/token"
import "strings"

func HasSymbol(source []byte, symbol string, declaration_type string) bool {

	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", source, 0)

	if err != nil {
		return false
	}

	if declaration_type == "func" {
		return hasFunction(file, fileset, symbol)
	} else if declaration_type == "interface" || declaration_type == "struct" {
		return hasType(file, symbol, declaration_type)
	}

	return false

}

func hasFunction(file *ast.File, fileset *token.FileSet, symbol string) bool {

	for _, declaration := range file.Decls {

		func_decl, ok := declaration.(*ast.FuncDecl)

		if ok == true && func_decl.Name != nil {

			if func_decl.Recv != nil && len(func_decl.Recv.List) > 0 {
				receiver_type := receiverTypeName(fileset, func_decl.Recv.List[0].Type)
				if receiver_type+"."+func_decl.Name.Name == symbol {
					return true
				}
			} else if func_decl.Name.Name == symbol {
				return true
			}

		}

	}

	return false

}

func hasType(file *ast.File, symbol string, declaration_type string) bool {

	for _, declaration := range file.Decls {

		gen_decl, ok := declaration.(*ast.GenDecl)

		if ok == true && gen_decl.Tok == token.TYPE {

			for _, specification := range gen_decl.Specs {

				type_spec, ok := specification.(*ast.TypeSpec)

				if ok == true && type_spec.Name != nil && type_spec.Name.Name == symbol {

					if declaration_type == "interface" {
						_, is_interface := type_spec.Type.(*ast.InterfaceType)
						if is_interface == true {
							return true
						}
					} else if declaration_type == "struct" {
						_, is_struct := type_spec.Type.(*ast.StructType)
						if is_struct == true {
							return true
						}
					}

				}

			}

		}

	}

	return false

}

func receiverTypeName(fileset *token.FileSet, recv_type ast.Expr) string {

	if star_expression, ok := recv_type.(*ast.StarExpr); ok == true {
		recv_type = star_expression.X
	}

	buffer := bytes.Buffer{}
	printer.Fprint(&buffer, fileset, recv_type)

	return strings.TrimSpace(buffer.String())

}
