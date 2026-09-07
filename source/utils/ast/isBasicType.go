package ast

import "go/ast"

func isBasicType(expr ast.Expr) bool {

	switch value := expr.(type) {

	case *ast.Ident:

		switch value.Name {
		case "bool",
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"uintptr",
			"byte", "rune",
			"float32", "float64",
			"complex64", "complex128",
			"string":
			return true
		default:
			return false
		}

	case *ast.ArrayType:
		return isBasicType(value.Elt)

	case *ast.MapType:
		return isBasicType(value.Key) && isBasicType(value.Value)

	default:
		return false

	}

}

