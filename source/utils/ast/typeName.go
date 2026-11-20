package ast

import "bytes"
import "go/ast"
import "go/printer"
import "go/token"
import "strings"

func typeName(fileset *token.FileSet, expr ast.Expr) string {

	if expr == nil {
		return ""
	} else if _, ok := expr.(*ast.FuncType); ok == true {
		return "func"
	} else if _, ok := expr.(*ast.InterfaceType); ok == true {
		return "interface"
	} else if _, ok := expr.(*ast.StructType); ok == true {
		return "struct"
	} else {

		buffer := bytes.Buffer{}
		printer.Fprint(&buffer, fileset, expr)

		return strings.TrimSpace(buffer.String())

	}
}
