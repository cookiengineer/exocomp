package lint

import "go/ast"
import "strconv"

func importPath(spec *ast.ImportSpec) string {

	if spec == nil || spec.Path == nil {
		return ""
	} else {

		value, err := strconv.Unquote(spec.Path.Value)

		if err == nil {
			return value
		} else {
			return ""
		}

	}

}
