package main

import "go/ast"

type methodReference struct {
	File *File
	Decl *ast.FuncDecl
	Name string
}
