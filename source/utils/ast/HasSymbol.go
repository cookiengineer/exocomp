package ast

import "go/parser"
import "go/token"

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

		return hasType(file, fileset, symbol, "func")

	} else if expected_type == "interface" {

		return hasType(file, fileset, symbol, "interface")

	} else if expected_type == "struct" {

		return hasType(file, fileset, symbol, "struct")

	} else {

		if hasFunc(file, fileset, symbol) == true {
			return true
		}

		return hasType(file, fileset, symbol, "")

	}

}
