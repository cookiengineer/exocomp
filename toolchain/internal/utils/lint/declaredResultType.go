package lint

import "go/ast"

// declaredResultType returns the declared package type a function returns.
func declaredResultType(func_decl *ast.FuncDecl, types map[string]*File) string {

	if func_decl.Type.Results == nil {
		return ""
	} else {

		for _, field := range func_decl.Type.Results.List {

			name := receiverTypeName(field.Type)

			if name != "" && types[name] != nil {
				return name
			}

		}

		return ""

	}

}
