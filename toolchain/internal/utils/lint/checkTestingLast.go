package lint

func checkTestingLast(_ *Context, file *File) []Diagnostic {

	if file.IsTest == false {
		return nil
	} else {

		entries := collectImports(file)
		last_stdlib := ""
		last_entry := importEntry{}

		for _, entry := range entries {

			if classifyImport(entry.Path) != groupStdlib {
				continue
			}

			last_stdlib = entry.Path
			last_entry = entry

		}

		if last_stdlib == "testing" {
			return nil
		} else {

			for _, entry := range entries {

				if entry.Path != "testing" {
					continue
				}

				return []Diagnostic{
					file.Diagnostic(
						"import/testing-last",
						last_entry.Spec.Pos(),
						"testing must be the last standard-library import",
					),
				}

			}

			return nil

		}

	}

}
