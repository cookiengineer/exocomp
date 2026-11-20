package lint

import "go/ast"

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
