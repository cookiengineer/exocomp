package ast

import "bytes"
import "go/ast"
import "go/printer"
import "go/token"
import "strings"

func getFuncReturnType(fileset *token.FileSet, func_decl *ast.FuncDecl) string {

	if func_decl.Type == nil || func_decl.Type.Results == nil || len(func_decl.Type.Results.List) == 0 {
		return ""
	}

	types := make([]string, 0)

	for _, field := range func_decl.Type.Results.List {

		buffer := bytes.Buffer{}
		printer.Fprint(&buffer, fileset, field.Type)
		type_name := strings.TrimSpace(buffer.String())

		if len(field.Names) > 0 {

			for range field.Names {
				types = append(types, type_name)
			}

		} else {
			types = append(types, type_name)
		}

	}

	if len(types) == 1 {
		return types[0]
	}

	return "(" + strings.Join(types, ", ") + ")"

}

