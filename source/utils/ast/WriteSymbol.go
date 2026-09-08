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

		result, _ = writeType(source, fileset, file, symbol, "func", replacement)

		return result

	} else if declaration_type == "interface" {

		result, _ := writeType(source, fileset, file, symbol, "interface", replacement)

		return result

	} else if declaration_type == "struct" {

		result, _ := writeType(source, fileset, file, symbol, "struct", replacement)

		return result

	} else {

		result, ok := writeFunc(source, fileset, file, symbol, replacement)

		if ok == true {
			return result
		}

		result, _ = writeType(source, fileset, file, symbol, "", replacement)

		return result

	}

}

