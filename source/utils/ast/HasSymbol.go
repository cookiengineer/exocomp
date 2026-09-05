package ast

import "bytes"
import "go/ast"
import "go/parser"
import "go/printer"
import "go/token"
import "strings"

func HasSymbol(source []byte, symbol string, expected_type string) bool {

	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", source, 0)

	if err != nil {
		return false
	}

	if expected_type == "func" {

		if hasFunc(file, fileset, symbol) == true {
			return true
		}

	}

	return hasType(file, fileset, symbol, expected_type)

}

func receiverTypeName(fileset *token.FileSet, recv_type ast.Expr) string {

	if star_expression, ok := recv_type.(*ast.StarExpr); ok == true {
		recv_type = star_expression.X
	}

	buffer := bytes.Buffer{}
	printer.Fprint(&buffer, fileset, recv_type)

	return strings.TrimSpace(buffer.String())

}
