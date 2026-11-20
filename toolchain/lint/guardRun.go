package main

import "go/ast"

type guardRun struct {
	Stmts   []ast.Stmt
	Start   int
	End     int
	HasTail bool
	Depth   int
}
