package lint

import "strings"

func checkTypesNoCapabilities(context *Context, file *File) []Diagnostic {

	if strings.HasPrefix(file.Rel, "types/") == false || context.Module == "" {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, imported := range importPaths(file) {

			if imported == context.Module+"/utils/capabilities" {

				result = append(
					result,
					file.Diagnostic(
						"dep/types-no-capabilities",
						file.AST.Pos(),
						"types must not import utils/capabilities",
					),
				)

			}

		}

		return result

	}

}
