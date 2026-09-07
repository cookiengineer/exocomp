package ast

import "go/ast"
import "go/token"

func writeFunc(source []byte, fileset *token.FileSet, file *ast.File, symbol string, replacement string) ([]byte, bool) {

	for _, decl := range file.Decls {

		func_decl, ok := decl.(*ast.FuncDecl)

		if ok == true && func_decl.Name != nil {

			if func_decl.Recv != nil && len(func_decl.Recv.List) > 0 {

				receiver_type := receiverTypeName(fileset, func_decl.Recv.List[0].Type)

				if receiver_type+"."+func_decl.Name.Name == symbol {
					return spliceNode(source, fileset, func_decl, replacement), true
				}

			} else if func_decl.Name.Name == symbol {
				return spliceNode(source, fileset, func_decl, replacement), true
			}

		}

	}

	return source, false

}

