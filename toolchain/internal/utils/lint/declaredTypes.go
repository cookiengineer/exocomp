package lint

import "go/ast"
import "go/token"

func declaredTypes(file *File) []*ast.TypeSpec {

	result := make([]*ast.TypeSpec, 0)

	for _, decl := range file.AST.Decls {

		gen_decl, ok := decl.(*ast.GenDecl)

		if ok == true && gen_decl.Tok == token.TYPE {

			for _, spec := range gen_decl.Specs {

				type_spec, ok := spec.(*ast.TypeSpec)

				if ok == true {
					result = append(result, type_spec)
				}

			}

		}

	}

	return result

}
