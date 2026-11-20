package main

import "go/ast"
import "go/token"
import "strings"

const (
	groupProject = 0
	groupStdlib  = 1
	groupThird   = 2
)

type importEntry struct {
	Path string
	Spec *ast.ImportSpec
	Line int
}

func classifyImport(path string) int {

	if strings.HasPrefix(path, "exocomp/") == true || strings.HasPrefix(path, "exocomp-toolchain/") == true {
		return groupProject
	}

	first := path

	if index := strings.Index(path, "/"); index >= 0 {
		first = path[:index]
	}

	if strings.Contains(first, ".") == true {
		return groupThird
	}

	return groupStdlib

}

func collectImports(file *File) []importEntry {

	result := make([]importEntry, 0, len(file.AST.Imports))

	for _, spec := range file.AST.Imports {
		result = append(result, importEntry{
			Path: importPath(spec),
			Spec: spec,
			Line: file.Line(spec.Pos()),
		})
	}

	return result

}

func checkImportParen(_ *Context, file *File) []Diagnostic {

	result := make([]Diagnostic, 0)

	for _, decl := range file.AST.Decls {

		gen_decl, ok := decl.(*ast.GenDecl)

		if ok == false || gen_decl.Tok != token.IMPORT {
			continue
		}

		if gen_decl.Lparen.IsValid() == true {
			result = append(
				result,
				file.Diagnostic(
					"import/paren",
					gen_decl.Pos(),
					"imports must be one per line, without a parenthesized block",
				),
			)
		}

	}

	return result

}

func checkImportGroup(_ *Context, file *File) []Diagnostic {

	entries := collectImports(file)

	if len(entries) == 0 {
		return nil
	}

	result := make([]Diagnostic, 0)

	previous := entries[0]
	seen_stdlib := false

	for index := 1; index < len(entries); index++ {

		entry := entries[index]
		group := classifyImport(entry.Path)
		previous_group := classifyImport(previous.Path)

		if group == groupThird && previous_group != groupThird {

			if entry.Line-previous.Line < 2 {
				result = append(
					result,
					file.Diagnostic(
						"import/group",
						entry.Spec.Pos(),
						"missing blank line before the third-party import group",
					),
				)
			}

		} else if group != groupThird && previous_group != groupThird {

			if entry.Line-previous.Line > 1 {
				result = append(
					result,
					file.Diagnostic(
						"import/group",
						entry.Spec.Pos(),
						"project and standard library imports share one block",
					),
				)
			}

		} else if group < previous_group {
			result = append(
				result,
				file.Diagnostic(
					"import/group",
					entry.Spec.Pos(),
					"import %q is in the %s group but comes after the %s group",
					entry.Path,
					groupName(group),
					groupName(previous_group),
				),
			)
		}

		if group == groupStdlib {
			seen_stdlib = true
		} else if group == groupProject && seen_stdlib == true {
			result = append(
				result,
				file.Diagnostic(
					"import/group",
					entry.Spec.Pos(),
					"project import %q must come before standard library imports",
					entry.Path,
				),
			)
		}

		previous = entry

	}

	return result

}

func groupName(group int) string {

	switch group {
	case groupProject:
		return "project"
	case groupStdlib:
		return "standard library"
	default:
		return "third-party"
	}

}

func checkTestingLast(_ *Context, file *File) []Diagnostic {

	if file.IsTest == false {
		return nil
	}

	entries := collectImports(file)
	last_stdlib := ""
	last_entry := importEntry{}

	for _, entry := range entries {

		if classifyImport(entry.Path) != groupStdlib {
			continue
		}

		last_stdlib = entry.Path
		last_entry = entry

	}

	if last_stdlib == "testing" {
		return nil
	}

	for _, entry := range entries {

		if entry.Path != "testing" {
			continue
		}

		return []Diagnostic{
			file.Diagnostic(
				"import/testing-last",
				last_entry.Spec.Pos(),
				"testing must be the last standard-library import",
			),
		}

	}

	return nil

}

func checkImportOrder(_ *Context, file *File) []Diagnostic {

	entries := collectImports(file)

	if len(entries) < 2 {
		return nil
	}

	result := make([]Diagnostic, 0)
	previous := entries[0]

	for _, entry := range entries[1:] {

		if classifyImport(entry.Path) == classifyImport(previous.Path) && previous.Path > entry.Path {
			result = append(
				result,
				file.Diagnostic(
					"import/order",
					entry.Spec.Pos(),
					"import %q must be sorted before %q",
					entry.Path,
					previous.Path,
				),
			)
		}

		previous = entry

	}

	return result

}
