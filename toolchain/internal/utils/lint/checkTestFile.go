package lint

import "go/ast"
import "strings"

func checkTestFile(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {

		sources := make(map[string]bool)

		for _, candidate := range file.Pkg.Files {

			if candidate.IsTest == false {
				sources[fileBaseName(candidate)] = true
			}

		}

		base := strings.TrimSuffix(fileBaseName(file), "_test")
		index := strings.LastIndex(base, "_")

		if sources[base] == true {
			return nil
		} else if index > 0 && sources[base[:index]] == true {
			return nil
		} else {

			has_test := false

			for _, decl := range file.AST.Decls {

				func_decl, ok := decl.(*ast.FuncDecl)

				if ok == true && strings.HasPrefix(func_decl.Name.Name, "Test") == true {
					has_test = true
				}

			}

			if has_test == false {
				return nil
			} else {
				return []Diagnostic{
					file.Diagnostic(
						"test/file",
						file.AST.Pos(),
						"test file %q has no matching source file",
						fileBaseName(file)+".go",
					),
				}
			}

		}

	} else {
		return nil
	}

}
