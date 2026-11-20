package main

import "fmt"
import "go/ast"
import "go/token"
import "strings"

// File is one parsed Go file.
type File struct {
	Path     string
	Rel      string
	Src      []byte
	AST      *ast.File
	Lines    []string
	Fset     *token.FileSet
	Pkg      *Package
	IsTest   bool
	Suppress map[int]map[string]bool
}

// Position returns the token position of pos.
func (file *File) Position(pos token.Pos) token.Position {
	return file.Fset.Position(pos)
}

// Line returns the 1-based line of pos.
func (file *File) Line(pos token.Pos) int {
	return file.Fset.Position(pos).Line
}

// Text returns the source slice covered by node.
func (file *File) Text(node ast.Node) string {

	start := file.Fset.Position(node.Pos()).Offset
	end := file.Fset.Position(node.End()).Offset

	if start < 0 || end > len(file.Src) || start > end {
		return ""
	} else {
		return strings.TrimSpace(string(file.Src[start:end]))
	}

}

// Diagnostic builds a diagnostic for rule at pos.
func (file *File) Diagnostic(rule string, pos token.Pos, format string, arguments ...any) Diagnostic {
	return Diagnostic{
		Position: file.Position(pos),
		Rule:     rule,
		Message:  fmt.Sprintf(format, arguments...),
	}
}

func (file *File) IsSuppressed(line int, rule string) bool {

	rules := file.Suppress[line]

	if rules == nil {
		return false
	} else if rules["all"] == true {
		return true
	} else {
		return rules[rule]
	}

}
