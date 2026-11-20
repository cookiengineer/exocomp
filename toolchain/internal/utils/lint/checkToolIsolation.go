package lint

import "strings"

func checkToolIsolation(context *Context, file *File) []Diagnostic {

	current, ok := isToolPackage(file.Rel)

	if ok == false || context.Module == "" {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, imported := range importPaths(file) {

			prefix := context.Module + "/tools/"

			if strings.HasPrefix(imported, prefix) == true {

				other := strings.Split(strings.TrimPrefix(imported, prefix), "/")[0]

				if other != current {

					result = append(
						result,
						file.Diagnostic(
							"dep/tool-isolation",
							file.AST.Pos(),
							"tool package %s must not import tool package %s",
							current,
							other,
						),
					)

				}

			}

		}

		return result

	}

}
