package lint

func checkImportGroup(_ *Context, file *File) []Diagnostic {

	entries := collectImports(file)

	if len(entries) == 0 {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		previous := entries[0]
		seen_stdlib := false

		for index := 1; index < len(entries); index++ {

			entry := entries[index]
			group := classifyImport(entry.Path)
			previous_group := classifyImport(previous.Path)

			if group == groupThird && previous_group != groupThird {

				if entry.Line-previous.Line < 2 {
					result = append(
						result,
						file.Diagnostic(
							"import/group",
							entry.Spec.Pos(),
							"missing blank line before the third-party import group",
						),
					)
				}

			} else if group != groupThird && previous_group != groupThird {

				if entry.Line-previous.Line > 1 {
					result = append(
						result,
						file.Diagnostic(
							"import/group",
							entry.Spec.Pos(),
							"project and standard library imports share one block",
						),
					)
				}

			} else if group < previous_group {
				result = append(
					result,
					file.Diagnostic(
						"import/group",
						entry.Spec.Pos(),
						"import %q is in the %s group but comes after the %s group",
						entry.Path,
						groupName(group),
						groupName(previous_group),
					),
				)
			}

			if group == groupStdlib {
				seen_stdlib = true
			} else if group == groupProject && seen_stdlib == true {
				result = append(
					result,
					file.Diagnostic(
						"import/group",
						entry.Spec.Pos(),
						"project import %q must come before standard library imports",
						entry.Path,
					),
				)
			}

			previous = entry

		}

		return result

	}

}
