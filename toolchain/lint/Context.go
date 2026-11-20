package main

import "go/token"

// Context is the parsed source tree.
type Context struct {
	Root        string
	Module      string
	Fset        *token.FileSet
	Files       []*File
	Packages    []*Package
	ParseErrors []string
}

// Run executes every rule and returns the sorted, filtered diagnostics.
func (context *Context) Run(selected map[string]bool) []Diagnostic {

	result := make([]Diagnostic, 0)

	for _, file := range context.Files {

		for _, rule := range fileRules {

			if selected != nil && selected[rule.ID] == false {
				continue
			}

			for _, diagnostic := range rule.Check(context, file) {

				diagnostic.Rule = rule.ID

				if file.IsSuppressed(diagnostic.Position.Line, rule.ID) == true {
					continue
				}

				result = append(result, diagnostic)

			}

		}

	}

	for _, pkg := range context.Packages {

		for _, rule := range packageRules {

			if selected != nil && selected[rule.ID] == false {
				continue
			}

			for _, diagnostic := range rule.Check(context, pkg) {

				diagnostic.Rule = rule.ID

				if context.isSuppressed(diagnostic) == true {
					continue
				}

				result = append(result, diagnostic)

			}

		}

	}

	sortDiagnostics(result)

	return result

}

func (context *Context) isSuppressed(diagnostic Diagnostic) bool {

	for _, file := range context.Files {

		if file.Path == diagnostic.Position.Filename {
			return file.IsSuppressed(diagnostic.Position.Line, diagnostic.Rule)
		}

	}

	return false

}
