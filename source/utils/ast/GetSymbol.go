package ast

import "go/parser"
import "go/token"

func GetSymbol(source []byte, symbol string, declaration_type string) string {

	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", source, 0)

	if err == nil {

		if declaration_type == "func" {

			result := getFunc(file, fileset, symbol)

			if result == "" {
				result = getType(file, fileset, symbol, declaration_type)
			}

			return result

		}

		return getType(file, fileset, symbol, declaration_type)

	}

	return ""

}
