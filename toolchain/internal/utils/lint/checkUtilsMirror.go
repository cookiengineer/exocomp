package lint

import "strings"

func checkUtilsMirror(context *Context, file *File) []Diagnostic {

	if strings.HasPrefix(file.Rel, "utils/") == false {
		return nil
	} else {

		parts := strings.Split(file.Rel, "/")

		if len(parts) < 3 {
			return nil
		} else {

			name := parts[1]

			for _, pkg := range context.Packages {

				if pkg.Rel == name && pkg.Name == name {
					return []Diagnostic{
						file.Diagnostic(
							"layout/utils-mirror",
							file.AST.Pos(),
							"utils/%s duplicates the top-level package %s",
							name,
							name,
						),
					}
				}

			}

			return nil

		}

	}

}
