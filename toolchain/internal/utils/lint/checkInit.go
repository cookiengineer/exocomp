package lint

import "go/ast"

func checkInit(_ *Context, file *File) []Diagnostic {

	count := 0

	for _, decl := range file.AST.Decls {

		func_decl, ok := decl.(*ast.FuncDecl)

		if ok == true && func_decl.Recv == nil && func_decl.Name.Name == "init" {
			count++
		}

	}

	if count <= 1 {
		return nil
	} else {
		return []Diagnostic{
			file.Diagnostic("layout/init", file.AST.Pos(), "file declares %d init() functions", count),
		}
	}

}
