package lint

func checkImportOrder(_ *Context, file *File) []Diagnostic {

	entries := collectImports(file)

	if len(entries) < 2 {
		return nil
	} else {

		result := make([]Diagnostic, 0)
		previous := entries[0]

		for _, entry := range entries[1:] {

			if classifyImport(entry.Path) == classifyImport(previous.Path) && previous.Path > entry.Path {
				result = append(
					result,
					file.Diagnostic(
						"import/order",
						entry.Spec.Pos(),
						"import %q must be sorted before %q",
						entry.Path,
						previous.Path,
					),
				)
			}

			previous = entry

		}

		return result

	}

}
