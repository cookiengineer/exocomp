package ast

import "go/parser"
import "go/token"

func GetSymbol(source []byte, symbol string, declaration_type string) *Symbol {

	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", source, 0)

	if err == nil {

		if declaration_type == "func" {

			result := getFunc(file, fileset, symbol)

			if result == nil {
				result = getType(file, fileset, symbol, declaration_type)
			}

			return result

		} else if declaration_type == "type" {
			return getType(file, fileset, symbol, declaration_type)
		}

	}

	return nil

}
