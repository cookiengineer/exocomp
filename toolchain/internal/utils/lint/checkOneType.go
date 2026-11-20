package lint

import "strings"

func checkOneType(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		names := make([]string, 0)

		for _, type_spec := range declaredTypes(file) {

			if isExported(type_spec.Name.Name) == true {
				names = append(names, type_spec.Name.Name)
			}

		}

		result := make([]Diagnostic, 0)

		for index, name := range names {

			if index > 0 {
				result = append(
					result,
					file.Diagnostic(
						"layout/one-type",
						file.AST.Pos(),
						"file declares %d exported types (%s); %s belongs in %s.go",
						len(names),
						strings.Join(names, ", "),
						name,
						name,
					),
				)
			}

		}

		return result

	}

}
