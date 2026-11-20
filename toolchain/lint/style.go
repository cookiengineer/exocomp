package main

import "go/ast"
import "go/format"
import "go/token"
import "regexp"
import "strconv"
import "strings"

var camelCasePattern = regexp.MustCompile(`^[a-z][a-z0-9]*([A-Z][a-z0-9]*)+$`)

var allowedLocalAcronyms = map[string]bool{
	"ANSI":  true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CVE":   true,
	"CVSS":  true,
	"CWE":   true,
	"DNS":   true,
	"EOF":   true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IO":    true,
	"IP":    true,
	"JSON":  true,
	"OS":    true,
	"PTY":   true,
	"SHA":   true,
	"SQL":   true,
	"TLS":   true,
	"TTL":   true,
	"TTY":   true,
	"UI":    true,
	"URI":   true,
	"URL":   true,
	"UUID":  true,
	"YAML":  true,
}

func checkAny(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	ast.Inspect(file.AST, func(node ast.Node) bool {

		interface_type, ok := node.(*ast.InterfaceType)

		if ok == false || interface_type.Methods == nil {
			return true
		}

		if len(interface_type.Methods.List) != 0 {
			return true
		}

		result = append(
			result,
			file.Diagnostic(
				"style/any",
				interface_type.Pos(),
				"use any instead of interface{}",
			),
		)

		return true

	})

	return result

}

func checkBlankBlock(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	ast.Inspect(file.AST, func(node ast.Node) bool {

		block, ok := node.(*ast.BlockStmt)

		if ok == false || len(block.List) <= 1 {
			return true
		}

		open_line := file.Line(block.Lbrace)
		close_line := file.Line(block.Rbrace)

		if open_line == close_line {
			return true
		}

		after := open_line + 1

		if after <= close_line && after-1 < len(file.Lines) {

			line := file.Lines[after-1]

			if strings.TrimSpace(line) != "" && isCommentLine(line) == false {
				result = append(
					result,
					file.Diagnostic(
						"style/blank-block",
						block.Lbrace,
						"multi-statement block needs a blank line after the opening brace",
					),
				)
			}

		}

		before := close_line - 1

		if before > open_line && before-1 < len(file.Lines) && before-1 >= 0 {

			line := file.Lines[before-1]

			if strings.TrimSpace(line) != "" && isCommentLine(line) == false {
				result = append(
					result,
					file.Diagnostic(
						"style/blank-block",
						block.Rbrace,
						"multi-statement block needs a blank line before the closing brace",
					),
				)
			}

		}

		return true

	})

	return result

}

func checkGofmt(_ *Context, file *File) []Diagnostic {

	formatted, err := format.Source(file.Src)

	if err != nil {
		return []Diagnostic{
			file.Diagnostic("style/gofmt", file.AST.Pos(), "cannot format: %s", err.Error()),
		}
	}

	if string(formatted) == string(file.Src) {
		return nil
	}

	return []Diagnostic{
		file.Diagnostic("style/gofmt", file.AST.Pos(), "file is not gofmt formatted"),
	}

}

func checkJSONIndent(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	ast.Inspect(file.AST, func(node ast.Node) bool {

		call, ok := node.(*ast.CallExpr)

		if ok == false {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)

		if ok == false || selector.Sel.Name != "MarshalIndent" {
			return true
		}

		ident, ok := selector.X.(*ast.Ident)

		if ok == false || ident.Name != "json" {
			return true
		}

		if len(call.Args) < 3 {
			return true
		}

		literal, ok := call.Args[2].(*ast.BasicLit)

		if ok == false || literal.Kind != token.STRING {
			return true
		}

		value, err := strconv.Unquote(literal.Value)

		if err != nil || value == "\t" {
			return true
		}

		result = append(
			result,
			file.Diagnostic(
				"style/json-indent",
				call.Pos(),
				"JSON indentation must be a tab, got %q",
				value,
			),
		)

		return true

	})

	return result

}

func isUpperRune(character rune) bool {
	return character >= 'A' && character <= 'Z'
}

func isLowerRune(character rune) bool {
	return character >= 'a' && character <= 'z'
}

func camelSegments(name string) []string {

	runes := []rune(name)
	segments := make([]string, 0)
	current := make([]rune, 0)

	for index, character := range runes {

		if isUpperRune(character) == true && len(current) > 0 {

			previous_upper := isUpperRune(runes[index-1])
			next_lower := index+1 < len(runes) && isLowerRune(runes[index+1])

			if previous_upper == false || next_lower == true {
				segments = append(segments, string(current))
				current = make([]rune, 0)
			}

		}

		current = append(current, character)

	}

	if len(current) > 0 {
		segments = append(segments, string(current))
	}

	return segments

}

func isCamelLocal(name string) bool {

	if camelCasePattern.MatchString(name) == false {
		return false
	}

	segments := camelSegments(name)

	for _, segment := range segments[1:] {

		is_upper := segment == strings.ToUpper(segment)

		if is_upper == true && allowedLocalAcronyms[segment] == true {
			continue
		}

		return true

	}

	return false

}

func checkLocalCase(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	}

	result := make([]Diagnostic, 0)

	report := func(name *ast.Ident) {

		if name == nil || name.Name == "_" {
			return
		}

		if isCamelLocal(name.Name) == false {
			return
		}

		result = append(
			result,
			file.Diagnostic(
				"style/local-case",
				name.Pos(),
				"local variable %s uses camelCase; use snake_case",
				name.Name,
			),
		)

	}

	reportList := func(fields *ast.FieldList) {

		if fields == nil {
			return
		}

		for _, field := range fields.List {
			for _, name := range field.Names {
				report(name)
			}
		}

	}

	for _, decl := range file.AST.Decls {

		func_decl, ok := decl.(*ast.FuncDecl)

		if ok == false || func_decl.Body == nil {
			continue
		}

		reportList(func_decl.Type.Params)

		ast.Inspect(func_decl.Body, func(node ast.Node) bool {

			switch typed := node.(type) {
			case *ast.FuncLit:
				reportList(typed.Type.Params)
			case *ast.AssignStmt:
				if typed.Tok != token.DEFINE {
					return true
				}
				for _, left := range typed.Lhs {
					ident, ok := left.(*ast.Ident)
					if ok == true {
						report(ident)
					}
				}
			case *ast.RangeStmt:
				if typed.Tok != token.DEFINE {
					return true
				}
				if ident, ok := typed.Key.(*ast.Ident); ok == true {
					report(ident)
				}
				if ident, ok := typed.Value.(*ast.Ident); ok == true {
					report(ident)
				}
			case *ast.ValueSpec:
				for _, name := range typed.Names {
					report(name)
				}
			}

			return true

		})

	}

	return result

}
