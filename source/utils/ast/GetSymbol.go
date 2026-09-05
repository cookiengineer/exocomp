package ast

import "go/parser"
import "go/token"

func GetSymbol(source []byte, symbol string, declaration_type string) string {

	fileset   := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", source, 0)

	if err == nil {

		if declaration_type == "func" {
			return getFunc(file, fileset, symbol)
		} else if declaration_type == "interface" {
			return getType(file, fileset, symbol, "interface")
		} else if declaration_type == "struct" {
			return getType(file, fileset, symbol, "struct")
		}

	}

	return ""

}
