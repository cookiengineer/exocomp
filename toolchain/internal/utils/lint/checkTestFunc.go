package lint

import "go/ast"
import "strings"

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
