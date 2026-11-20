package main

import "go/ast"
import "go/token"
import "strings"

// statementList is one slice of statements plus its control-flow depth.
type statementList struct {
	Stmts    []ast.Stmt
	Depth    int
	Receiver string
}

// collectStatementLists gathers every statement list in the file with its
// nesting depth. A function body has depth 0.
func collectStatementLists(file *File) []statementList {

	result := make([]statementList, 0)

	var walk_statements func([]ast.Stmt, int, string)
	var walk_statement func(ast.Stmt, int, string)

	walk_statements = func(stmts []ast.Stmt, depth int, receiver string) {

		if len(stmts) == 0 {
			return
		}

		result = append(result, statementList{Stmts: stmts, Depth: depth, Receiver: receiver})

		for _, stmt := range stmts {
			walk_statement(stmt, depth, receiver)
		}

	}

	walk_statement = func(stmt ast.Stmt, depth int, receiver string) {

		switch typed := stmt.(type) {
		case *ast.BlockStmt:
			walk_statements(typed.List, depth+1, receiver)
		case *ast.IfStmt:
			walk_statements(typed.Body.List, depth+1, receiver)

			switch alternative := typed.Else.(type) {
			case *ast.IfStmt:
				walk_statement(alternative, depth, receiver)
			case *ast.BlockStmt:
				walk_statements(alternative.List, depth+1, receiver)
			}
		case *ast.ForStmt:
			walk_statements(typed.Body.List, depth, receiver)
		case *ast.RangeStmt:
			walk_statements(typed.Body.List, depth, receiver)
		case *ast.SwitchStmt:
			for _, clause := range typed.Body.List {
				case_clause, ok := clause.(*ast.CaseClause)
				if ok == true {
					walk_statements(case_clause.Body, depth+1, receiver)
				}
			}
		case *ast.TypeSwitchStmt:
			for _, clause := range typed.Body.List {
				case_clause, ok := clause.(*ast.CaseClause)
				if ok == true {
					walk_statements(case_clause.Body, depth+1, receiver)
				}
			}
		case *ast.SelectStmt:
			for _, clause := range typed.Body.List {
				comm_clause, ok := clause.(*ast.CommClause)
				if ok == true {
					walk_statements(comm_clause.Body, depth+1, receiver)
				}
			}
		case *ast.LabeledStmt:
			walk_statement(typed.Stmt, depth, receiver)
		case *ast.CaseClause:
			walk_statements(typed.Body, depth, receiver)
		case *ast.CommClause:
			walk_statements(typed.Body, depth, receiver)
		}

	}

	for _, decl := range file.AST.Decls {

		func_decl, ok := decl.(*ast.FuncDecl)

		if ok == true && func_decl.Body != nil {

			receiver := ""

			if func_decl.Recv != nil && len(func_decl.Recv.List) > 0 {
				receiver = receiverName(func_decl.Recv.List[0])
			}

			walk_statements(func_decl.Body.List, 0, receiver)
		}

	}

	ast.Inspect(file.AST, func(node ast.Node) bool {

		function_literal, ok := node.(*ast.FuncLit)

		if ok == true {
			walk_statements(function_literal.Body.List, 0, "")
		}

		return true

	})

	return result

}

// guardKind classifies a block that always leaves the enclosing flow.
// It returns "" when the block is not a terminal guard.
func guardKind(block *ast.BlockStmt) string {

	if len(block.List) == 0 {
		return ""
	}

	kind := ""

	for _, stmt := range block.List {

		switch typed := stmt.(type) {
		case *ast.ReturnStmt:
			if len(typed.Results) > 0 {
				if kind == "" || kind == "bare" {
					kind = "value"
				}
			} else if kind == "" {
				kind = "bare"
			}
		case *ast.BranchStmt:
			if typed.Tok == token.FALLTHROUGH {
				return ""
			}
			if kind == "" {
				kind = "loop"
			}
		case *ast.ExprStmt:

			call, ok := typed.X.(*ast.CallExpr)

			if ok == false {
				return ""
			}

			if isFatalCall(call) == false {
				return ""
			}

			if kind == "" {
				kind = "fatal"
			}
		default:
			return ""
		}

	}

	return kind

}

func isFatalCall(call *ast.CallExpr) bool {

	switch typed := call.Fun.(type) {
	case *ast.Ident:
		return typed.Name == "panic"
	case *ast.SelectorExpr:

		name := typed.Sel.Name

		if name == "Exit" {
			return true
		}

		return strings.HasPrefix(name, "Fatal") == true || strings.HasPrefix(name, "Skip") == true
	default:
		return false
	}

}

func isValueGuard(stmt ast.Stmt) bool {

	if_stmt, ok := stmt.(*ast.IfStmt)

	if ok == false || if_stmt.Else != nil {
		return false
	}

	return guardKind(if_stmt.Body) == "value"

}

type guardRun struct {
	Stmts   []ast.Stmt
	Start   int
	End     int
	HasTail bool
	Depth   int
}

func collectGuardRuns(file *File) []guardRun {

	result := make([]guardRun, 0)

	for _, list := range collectStatementLists(file) {

		index := 0

		if list.Depth >= 3 {
			continue
		}

		for index < len(list.Stmts) {

			if isValueGuard(list.Stmts[index]) == false || isReceiverNilGuard(list.Stmts[index], list.Receiver) == true {
				index++
				continue
			}

			end := index

			for end < len(list.Stmts) && isValueGuard(list.Stmts[end]) == true {
				end++
			}

			result = append(result, guardRun{
				Stmts:   list.Stmts,
				Start:   index,
				End:     end,
				HasTail: end < len(list.Stmts),
				Depth:   list.Depth,
			})

			index = end

		}

	}

	return result

}

func checkGuard(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	}

	result := make([]Diagnostic, 0)

	for _, run := range collectGuardRuns(file) {

		if run.HasTail == false || run.End-run.Start != 1 {
			continue
		}

		result = append(
			result,
			file.Diagnostic(
				"branch/guard",
				run.Stmts[run.Start].Pos(),
				"guard clause has %d statements after it; invert it into a positive branch",
				len(run.Stmts)-1-run.Start,
			),
		)

	}

	return result

}

func checkGuardChain(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	}

	result := make([]Diagnostic, 0)

	for _, run := range collectGuardRuns(file) {

		if run.HasTail == false || run.End-run.Start < 2 {
			continue
		}

		result = append(
			result,
			file.Diagnostic(
				"branch/guard-chain",
				run.Stmts[run.Start].Pos(),
				"%d consecutive guards; use one } else if { chain",
				run.End-run.Start,
			),
		)

	}

	return result

}

// isReceiverNilGuard reports whether stmt is a `receiver == nil` guard. The
// codestyle permits early returns for nil receivers.
func isReceiverNilGuard(stmt ast.Stmt, receiver string) bool {

	if receiver == "" {
		return false
	}

	if_stmt, ok := stmt.(*ast.IfStmt)

	if ok == false {
		return false
	}

	binary, ok := unwrapParens(if_stmt.Cond).(*ast.BinaryExpr)

	if ok == false || (binary.Op != token.EQL && binary.Op != token.NEQ) {
		return false
	}

	left, left_ok := unwrapParens(binary.X).(*ast.Ident)
	right, right_ok := unwrapParens(binary.Y).(*ast.Ident)

	if left_ok == false || right_ok == false {
		return false
	}

	if left.Name == receiver && right.Name == "nil" {
		return true
	}

	return right.Name == receiver && left.Name == "nil"

}
