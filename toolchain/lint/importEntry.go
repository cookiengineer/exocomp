package main

import "go/ast"

type importEntry struct {
	Path string
	Spec *ast.ImportSpec
	Line int
}
