package ast

import "go/ast"
import "go/token"

func getSymbols(file *ast.File, fileset *token.FileSet) map[string]string {

	result := make(map[string]string, 0)

	for _, decl := range file.Decls {

		func_decl, ok1 := decl.(*ast.FuncDecl)

		if ok1 == true && func_decl.Name != nil {

			if func_decl.Recv != nil && len(func_decl.Recv.List) > 0 {

				receiver_type := receiverTypeName(fileset, func_decl.Recv.List[0].Type)
				result[receiver_type+"."+func_decl.Name.Name] = "func"

			} else {
				result[func_decl.Name.Name] = "func"
			}

		} else {

			gen_decl, ok2 := decl.(*ast.GenDecl)

			if ok2 == true && gen_decl.Tok == token.TYPE {

				for _, spec := range gen_decl.Specs {

					type_spec, ok3 := spec.(*ast.TypeSpec)

					if ok3 == true && type_spec.Name != nil {
						result[type_spec.Name.Name] = typeName(fileset, type_spec.Type)
					}

				}

			}

		}

	}

	return result

}
