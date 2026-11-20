package lint

import "strings"

func checkUtilsNoTools(context *Context, file *File) []Diagnostic {

	if strings.HasPrefix(file.Rel, "utils/") == false || context.Module == "" {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, imported := range importPaths(file) {

			if strings.HasPrefix(imported, context.Module+"/tools") == true {

				result = append(
					result,
					file.Diagnostic(
						"dep/utils-no-tools",
						file.AST.Pos(),
						"utils must not import %q",
						imported,
					),
				)

			}

		}

		return result

	}

}
