package lint

import "go/ast"

func checkMethodFile(_ *Context, pkg *Package) []Diagnostic {

	result := make([]Diagnostic, 0)

	for _, file := range pkg.Files {

		for _, decl := range file.AST.Decls {

			func_decl, ok := decl.(*ast.FuncDecl)

			if ok == false || func_decl.Recv == nil || len(func_decl.Recv.List) == 0 {
				continue
			}

			type_name := receiverTypeName(func_decl.Recv.List[0].Type)

			if type_name == "" {
				continue
			}

			owner := pkg.Types[type_name]

			if owner == nil || owner == file {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"layout/method-file",
					func_decl.Pos(),
					"method (%s *%s).%s must live in %s.go",
					receiverName(func_decl.Recv.List[0]),
					type_name,
					func_decl.Name.Name,
					type_name,
				),
			)

		}

	}

	return result

}
