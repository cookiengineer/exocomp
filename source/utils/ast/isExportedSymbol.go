package ast

import "strings"

func isExportedSymbol(name string) bool {
	return strings.ToUpper(name[0:1]) == name[0:1]
}
