package lint

func checkConstructorFile(_ *Context, pkg *Package) []Diagnostic {

	result := make([]Diagnostic, 0)

	for _, file := range pkg.Files {

		if file.IsTest == true {
			continue
		}

		for _, func_decl := range freeFunctions(file) {

			type_name := ownedTypeName(func_decl, pkg.Types)

			if type_name == "" {
				continue
			}

			if identifierBaseName(file) == type_name {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"layout/constructor-file",
					func_decl.Pos(),
					"%s builds %s and must live in %s.go",
					func_decl.Name.Name,
					type_name,
					type_name,
				),
			)

		}

	}

	return result

}
