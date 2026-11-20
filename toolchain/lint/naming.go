package main

import "go/ast"
import "path/filepath"
import "regexp"
import "strings"

var initialismPattern = regexp.MustCompile(`^[A-Z]{2,5}$`)

var allowedInitialisms = map[string]bool{
	"ANSI":  true,
	"API":   true,
	"ASCII": true,
	"CIDR":  true,
	"CPU":   true,
	"HTML":  true,
	"HTTP":  true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"OS":    true,
	"TTY":   true,
	"UI":    true,
	"URL":   true,
	"YAML":  true,
}

var predicatePrefixes = []string{
	"allow_", "can_", "did_", "has_", "in_", "is_", "needs_", "should_", "use_", "was_", "with_",
}

var predicateExceptions = map[string]bool{
	"_":     true,
	"ok":    true,
	"value": true,
}

func checkFileName(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		base := fileBaseName(file)

		if strings.Contains(base, "_") == false {
			return nil
		} else if hasPlatformSuffix(base) == true {
			return nil
		} else {
			return []Diagnostic{
				file.Diagnostic(
					"naming/file",
					file.AST.Pos(),
					"file name %q uses underscores; name it after the identifier it contains",
					filepath.Base(file.Rel),
				),
			}
		}

	}

}

func checkFuncNaming(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, decl := range file.AST.Decls {

			func_decl, ok := decl.(*ast.FuncDecl)

			if ok == false {
				continue
			}

			name := func_decl.Name.Name

			if isExported(name) == true || strings.Contains(name, "_") == false {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"naming/func",
					func_decl.Pos(),
					"unexported function %s uses snake_case; use camelCase",
					name,
				),
			)

		}

		return result

	}

}

func checkFromNaming(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, func_decl := range freeFunctions(file) {

			name := func_decl.Name.Name

			if strings.HasPrefix(name, "From") == false || len(name) <= len("From") {
				continue
			}

			next := name[len("From")]

			if next < 'A' || next > 'Z' {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"naming/from",
					func_decl.Pos(),
					"converter %s must be named <Type>From<Source>",
					name,
				),
			)

		}

		return result

	}

}

func checkTypeInitialism(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	for _, type_spec := range declaredTypes(file) {

		name := type_spec.Name.Name

		if initialismPattern.MatchString(name) == false {
			continue
		}

		if allowedInitialisms[name] == true {
			continue
		}

		result = append(
			result,
			file.Diagnostic(
				"naming/type-initialism",
				type_spec.Pos(),
				"type name %s is collapsed initials; spell it out",
				name,
			),
		)

	}

	return result

}

func isPredicateName(name string) bool {

	if predicateExceptions[name] == true {
		return false
	} else {

		for _, prefix := range predicatePrefixes {

			if strings.HasPrefix(name, prefix) == true {
				return false
			}

		}

		return true

	}

}

func isBoolType(expr ast.Expr) bool {

	if expr == nil {
		return false
	} else {

		ident, ok := expr.(*ast.Ident)

		return ok == true && ident.Name == "bool"

	}

}

func checkBooleanNaming(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		report := func(name *ast.Ident) {

			if name == nil || name.Name == "_" {
				return
			}

			if isPredicateName(name.Name) == false {
				return
			}

			result = append(
				result,
				file.Diagnostic(
					"naming/boolean",
					name.Pos(),
					"boolean %s does not read as a predicate; use is_/has_/can_/should_",
					name.Name,
				),
			)

		}

		ast.Inspect(file.AST, func(node ast.Node) bool {

			switch typed := node.(type) {

			case *ast.FuncDecl:

				if typed.Type.Params != nil {

					for _, field := range typed.Type.Params.List {

						if isBoolType(field.Type) == false {
							continue
						}

						for _, name := range field.Names {
							report(name)
						}

					}

				}

			case *ast.ValueSpec:

				if isBoolType(typed.Type) == true {

					for _, name := range typed.Names {
						report(name)
					}

				}

			}

			return true

		})

		return result

	}

}

func checkReceiverNaming(_ *Context, pkg *Package) []Diagnostic {

	by_type := make(map[string][]methodReference)

	for _, file := range pkg.Files {

		for _, decl := range file.AST.Decls {

			func_decl, ok := decl.(*ast.FuncDecl)

			if ok == false || func_decl.Recv == nil || len(func_decl.Recv.List) == 0 {
				continue
			}

			type_name := receiverTypeName(func_decl.Recv.List[0].Type)
			name := receiverName(func_decl.Recv.List[0])

			if type_name == "" || name == "" {
				continue
			}

			by_type[type_name] = append(by_type[type_name], methodReference{
				File: file,
				Decl: func_decl,
				Name: name,
			})

		}

	}

	result := make([]Diagnostic, 0)

	for type_name, references := range by_type {

		counts := make(map[string]int)

		for _, reference := range references {
			counts[reference.Name]++
		}

		if len(counts) <= 1 {
			continue
		}

		expected := ""

		for name := range counts {

			if expected == "" || counts[name] > counts[expected] {
				expected = name
			}

		}

		for _, reference := range references {

			if reference.Name == expected {
				continue
			}

			result = append(
				result,
				reference.File.Diagnostic(
					"naming/receiver",
					reference.Decl.Pos(),
					"receiver %s differs from %s used by the other %s methods",
					reference.Name,
					expected,
					type_name,
				),
			)

		}

	}

	return result

}
