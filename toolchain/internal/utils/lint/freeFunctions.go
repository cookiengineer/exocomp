package lint

import "go/ast"

func freeFunctions(file *File) []*ast.FuncDecl {

	result := make([]*ast.FuncDecl, 0)

	for _, decl := range file.AST.Decls {

		func_decl, ok := decl.(*ast.FuncDecl)

		if ok == true && func_decl.Recv == nil && func_decl.Name.Name != "init" {
			result = append(result, func_decl)
		}

	}

	return result

}
