package ast

import "fmt"
import "go/ast"
import "go/parser"
import "go/token"

func GetSymbolType(source []byte, symbol string) (string, error) {

	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", source, 0)

	if err != nil {
		return "", fmt.Errorf("Invalid Go syntax. \"type %s <type>\" must be defined!", symbol)
	}

	for _, decl := range file.Decls {

		gen_decl,  ok1 := decl.(*ast.GenDecl)
		func_decl, ok2 := decl.(*ast.FuncDecl)

		if ok1 == true && gen_decl.Tok == token.TYPE {

			for _, spec := range gen_decl.Specs {

				type_spec, ok3 := spec.(*ast.TypeSpec)

				if ok3 == true && type_spec.Name != nil && type_spec.Name.Name == symbol {

					if isBasicType(type_spec.Type) == false {
						return "", fmt.Errorf("Invalid type for symbol \"%s\". Only basic types and composites of basic types are supported.", symbol)
					}

					return typeName(fileset, type_spec.Type), nil

				}

			}

		} else if ok2 == true && func_decl.Name != nil {

			if func_decl.Recv != nil && len(func_decl.Recv.List) > 0 {

				receiver_type := receiverTypeName(fileset, func_decl.Recv.List[0].Type)

				if receiver_type+"."+func_decl.Name.Name == symbol {
					return getFuncReturnType(fileset, func_decl), nil
				}

			} else if func_decl.Name.Name == symbol {
				return getFuncReturnType(fileset, func_decl), nil
			}

		}

	}

	return "", fmt.Errorf("Invalid Go syntax. \"type %s <type>\" must be defined!", symbol)

}

