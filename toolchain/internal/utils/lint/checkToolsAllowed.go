package lint

import "strings"

func checkToolsAllowed(context *Context, file *File) []Diagnostic {

	if _, ok := isToolPackage(file.Rel); ok == false || context.Module == "" {
		return nil
	} else {

		prefix := context.Module + "/"

		allowed := func(imported string) bool {

			if strings.HasPrefix(imported, prefix) == false {
				return true
			} else {

				rest := strings.TrimPrefix(imported, prefix)

				if rest == "types" || rest == "schemas" || rest == "agents" {
					return true
				} else {
					return strings.HasPrefix(rest, "utils/") == true || strings.HasPrefix(rest, "parsers/") == true
				}

			}

		}

		result := make([]Diagnostic, 0)

		for _, imported := range importPaths(file) {

			if allowed(imported) == false {

				result = append(
					result,
					file.Diagnostic(
						"dep/tools-allowed",
						file.AST.Pos(),
						"tool packages may only import types, schemas, agents and utils, got %q",
						imported,
					),
				)

			}

		}

		return result

	}

}
