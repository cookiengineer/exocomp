package ast

import "bytes"
import "go/ast"
import "go/printer"
import "go/token"
import "strings"

func receiverTypeName(fileset *token.FileSet, recv_type ast.Expr) string {

	if star_expression, ok := recv_type.(*ast.StarExpr); ok == true {
		recv_type = star_expression.X
	}

	buffer := bytes.Buffer{}
	printer.Fprint(&buffer, fileset, recv_type)

	return strings.TrimSpace(buffer.String())

}
