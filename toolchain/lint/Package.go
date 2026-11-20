package main

import "go/token"

// Package is one directory of Go files.
type Package struct {
	Dir   string
	Rel   string
	Name  string
	Files []*File
	Types map[string]*File
}

// Diagnostic builds a diagnostic anchored to a file position.
func (pkg *Package) Diagnostic(rule string, file *File, pos token.Pos, format string, arguments ...any) Diagnostic {
	return file.Diagnostic(rule, pos, format, arguments...)
}
