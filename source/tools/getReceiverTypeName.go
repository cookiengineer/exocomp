package tools

import "bytes"
import "go/ast"
import "go/printer"
import "go/token"
import "strings"

func getReceiverTypeName(recv_type ast.Expr) string {

	if star_expression, ok := recv_type.(*ast.StarExpr); ok == true {
		recv_type = star_expression.X
	}

	buffer := bytes.Buffer{}
	printer.Fprint(&buffer, token.NewFileSet(), recv_type)

	return strings.TrimSpace(buffer.String())

}
