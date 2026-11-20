package main

import "go/ast"
import "regexp"
import "strings"

var testFunctionPattern = regexp.MustCompile(`^Test[A-Z][A-Za-z0-9]*(_[A-Za-z0-9]+)*$`)

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

func checkTestFunc(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {

		result := make([]Diagnostic, 0)

		for _, decl := range file.AST.Decls {

			func_decl, ok := decl.(*ast.FuncDecl)

			if ok == false || func_decl.Recv != nil {
				continue
			}

			name := func_decl.Name.Name

			if strings.HasPrefix(name, "Test") == false {
				continue
			}

			if name == "TestMain" {
				continue
			}

			if testFunctionPattern.MatchString(name) == true {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"test/func",
					func_decl.Pos(),
					"test function %s must be named Test<Type>_<Method>",
					name,
				),
			)

		}

		return result

	} else {
		return nil
	}

}
