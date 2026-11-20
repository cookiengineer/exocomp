package main

import "go/ast"

// statementList is one slice of statements plus its control-flow depth.
type statementList struct {
	Stmts    []ast.Stmt
	Depth    int
	Receiver string
}
