package lint

import "go/ast"

func checkReceiverNaming(_ *Context, pkg *Package) []Diagnostic {

	by_type := make(map[string][]methodReference)

	for _, file := range pkg.Files {

		for _, decl := range file.AST.Decls {

			func_decl, ok := decl.(*ast.FuncDecl)

			if ok == false || func_decl.Recv == nil || len(func_decl.Recv.List) == 0 {
				continue
			}

			type_name := receiverTypeName(func_decl.Recv.List[0].Type)
			name := receiverName(func_decl.Recv.List[0])

			if type_name == "" || name == "" {
				continue
			}

			by_type[type_name] = append(by_type[type_name], methodReference{
				File: file,
				Decl: func_decl,
				Name: name,
			})

		}

	}

	result := make([]Diagnostic, 0)

	for type_name, references := range by_type {

		counts := make(map[string]int)

		for _, reference := range references {
			counts[reference.Name]++
		}

		if len(counts) <= 1 {
			continue
		}

		expected := ""

		for name := range counts {

			if expected == "" || counts[name] > counts[expected] {
				expected = name
			}

		}

		for _, reference := range references {

			if reference.Name == expected {
				continue
			}

			result = append(
				result,
				reference.File.Diagnostic(
					"naming/receiver",
					reference.Decl.Pos(),
					"receiver %s differs from %s used by the other %s methods",
					reference.Name,
					expected,
					type_name,
				),
			)

		}

	}

	return result

}
