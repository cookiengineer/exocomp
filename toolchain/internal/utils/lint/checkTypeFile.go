package lint

func checkTypeFile(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		types := declaredTypes(file)

		if len(types) == 0 {
			return nil
		} else {

			has_exported := false

			for _, type_spec := range types {

				if isExported(type_spec.Name.Name) == true {
					has_exported = true
				}

			}

			if has_exported == true {
				return nil
			} else {

				base := identifierBaseName(file)
				result := make([]Diagnostic, 0)

				for _, type_spec := range types {

					if type_spec.Name.Name != base {
						result = append(
							result,
							file.Diagnostic(
								"layout/type-file",
								type_spec.Pos(),
								"unexported type %s must live in %s.go",
								type_spec.Name.Name,
								type_spec.Name.Name,
							),
						)
					}

				}

				return result

			}

		}

	}

}
