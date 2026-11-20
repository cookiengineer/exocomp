package lint

import "go/ast"

func receiverName(field *ast.Field) string {

	if len(field.Names) == 0 {
		return ""
	} else {
		return field.Names[0].Name
	}

}
