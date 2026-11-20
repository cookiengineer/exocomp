package main

import "go/ast"
import "go/parser"
import "go/token"
import "io/fs"
import "os"
import "path/filepath"
import "sort"
import "strings"

// Context is the parsed source tree.
type Context struct {
	Root        string
	Module      string
	Fset        *token.FileSet
	Files       []*File
	Packages    []*Package
	ParseErrors []string
}

// Package is one directory of Go files.
type Package struct {
	Dir   string
	Rel   string
	Name  string
	Files []*File
	Types map[string]*File
}

// File is one parsed Go file.
type File struct {
	Path     string
	Rel      string
	Src      []byte
	AST      *ast.File
	Lines    []string
	Fset     *token.FileSet
	Pkg      *Package
	IsTest   bool
	Suppress map[int]map[string]bool
}

var skippedFolders = map[string]bool{
	".git":         true,
	"node_modules": true,
	"testdata":     true,
	"vendor":       true,
}

// Load parses every Go file below root.
func Load(root string, include_tests bool) (*Context, error) {

	absolute, err0 := filepath.Abs(root)

	if err0 != nil {
		return nil, err0
	}

	context := &Context{
		Root:   absolute,
		Fset:   token.NewFileSet(),
		Module: readModuleName(absolute),
	}

	paths := make([]string, 0)

	err1 := filepath.WalkDir(absolute, func(path string, entry fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if entry.IsDir() == true {

			if path != absolute && skippedFolders[entry.Name()] == true {
				return fs.SkipDir
			}

			return nil

		}

		if strings.HasSuffix(path, ".go") == false {
			return nil
		}

		if include_tests == false && strings.HasSuffix(path, "_test.go") == true {
			return nil
		}

		paths = append(paths, path)

		return nil

	})

	if err1 != nil {
		return nil, err1
	}

	sort.Strings(paths)

	by_dir := make(map[string]*Package)

	for _, path := range paths {

		file, err2 := loadFile(context, absolute, path)

		if err2 != nil {
			context.ParseErrors = append(context.ParseErrors, err2.Error())
			continue
		}

		if file == nil {
			continue
		}

		directory := filepath.Dir(path)

		if by_dir[directory] == nil {
			relative, _ := filepath.Rel(absolute, directory)
			by_dir[directory] = &Package{
				Dir:   directory,
				Rel:   filepath.ToSlash(relative),
				Name:  file.AST.Name.Name,
				Types: make(map[string]*File),
			}
		}

		file.Pkg = by_dir[directory]
		by_dir[directory].Files = append(by_dir[directory].Files, file)
		context.Files = append(context.Files, file)

	}

	directories := make([]string, 0, len(by_dir))

	for directory := range by_dir {
		directories = append(directories, directory)
	}

	sort.Strings(directories)

	for _, directory := range directories {

		pkg := by_dir[directory]

		for _, file := range pkg.Files {

			for _, decl := range file.AST.Decls {

				gen_decl, ok := decl.(*ast.GenDecl)

				if ok == true && gen_decl.Tok == token.TYPE {

					for _, spec := range gen_decl.Specs {

						type_spec, ok := spec.(*ast.TypeSpec)

						if ok == true {
							pkg.Types[type_spec.Name.Name] = file
						}

					}

				}

			}

		}

		context.Packages = append(context.Packages, pkg)

	}

	return context, nil

}

func loadFile(context *Context, root string, path string) (*File, error) {

	source, err0 := os.ReadFile(path)

	if err0 != nil {
		return nil, err0
	}

	parsed, err1 := parser.ParseFile(context.Fset, path, source, parser.ParseComments|parser.AllErrors)

	if err1 != nil {
		return nil, err1
	}

	relative, err2 := filepath.Rel(root, path)

	if err2 != nil {
		relative = path
	}

	file := &File{
		Path:   path,
		Rel:    filepath.ToSlash(relative),
		Src:    source,
		AST:    parsed,
		Lines:  strings.Split(string(source), "\n"),
		Fset:   context.Fset,
		IsTest: strings.HasSuffix(path, "_test.go"),
	}

	file.Suppress = parseSuppressions(context.Fset, parsed)

	return file, nil

}

func readModuleName(root string) string {

	source, err0 := os.ReadFile(filepath.Join(root, "go.mod"))

	if err0 != nil {
		return ""
	}

	for _, line := range strings.Split(string(source), "\n") {

		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "module ") == true {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "module "))
		}

	}

	return ""

}

func isCommentLine(line string) bool {

	trimmed := strings.TrimSpace(line)

	return strings.HasPrefix(trimmed, "//") == true || strings.HasPrefix(trimmed, "/*") == true

}
