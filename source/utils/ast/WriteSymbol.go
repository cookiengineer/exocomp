package ast

import "go/parser"
import "go/token"

func WriteSymbol(source []byte, symbol string, declaration string, declaration_type string) []byte {

	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", source, 0)

	if err != nil {
		return source
	}

	replacement, name := parseDeclaration(declaration, declaration_type)

	if replacement == "" || name != symbol {
		return source
	}

	if declaration_type == "func" {

		result, ok := writeFunc(source, fileset, file, symbol, replacement)

		if ok == true {
			return result
		}

		result, _ = writeType(source, fileset, file, symbol, declaration_type, replacement)

		return result

	}

	result, _ := writeType(source, fileset, file, symbol, declaration_type, replacement)

	return result

}

