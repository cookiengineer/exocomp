package main

import "go/ast"
import "go/token"
import "path"
import "strings"

// oneFunctionExempt lists files that group low-level helpers on purpose. They
// are documented exceptions in docs/CODESTYLE.md.
var oneFunctionExempt = map[string]bool{
	"utils/landlock/landlock_linux.go": true,
}

var platformParts = map[string]bool{
	"aix": true, "android": true, "darwin": true, "dragonfly": true,
	"freebsd": true, "illumos": true, "ios": true, "js": true,
	"linux": true, "netbsd": true, "openbsd": true, "other": true,
	"plan9": true, "solaris": true, "unix": true, "wasip1": true,
	"windows": true,
	"386":     true, "amd64": true, "arm": true, "arm64": true,
	"loong64": true, "mips": true, "mips64": true, "mips64le": true,
	"mipsle": true, "ppc64": true, "ppc64le": true, "riscv64": true,
	"s390x": true, "wasm": true,
}

func fileBaseName(file *File) string {
	return strings.TrimSuffix(path.Base(file.Rel), ".go")
}

// identifierBaseName is the file name without platform suffixes.
func identifierBaseName(file *File) string {

	base := fileBaseName(file)
	parts := strings.Split(base, "_")

	for len(parts) > 1 {

		last := parts[len(parts)-1]

		if platformParts[last] == true {
			parts = parts[:len(parts)-1]
		} else {
			break
		}

	}

	return strings.Join(parts, "_")

}

// declaredResultType returns the declared package type a function returns.
func declaredResultType(func_decl *ast.FuncDecl, types map[string]*File) string {

	if func_decl.Type.Results == nil {
		return ""
	}

	for _, field := range func_decl.Type.Results.List {

		name := receiverTypeName(field.Type)

		if name != "" && types[name] != nil {
			return name
		}

	}

	return ""

}

// ownedTypeName returns the type a free function belongs to, if any.
// Constructors, parsers and converters live in their type's file.
func ownedTypeName(func_decl *ast.FuncDecl, types map[string]*File) string {

	name := func_decl.Name.Name

	for _, prefix := range []string{"New", "Parse", "new"} {

		if strings.HasPrefix(name, prefix) == false {
			continue
		}

		rest := strings.TrimPrefix(name, prefix)
		longest := ""

		for type_name := range types {

			if rest == type_name || strings.EqualFold(rest, type_name) == true || strings.HasPrefix(rest, type_name) == true {

				if len(type_name) > len(longest) {
					longest = type_name
				}

			}

		}

		if longest != "" {
			return longest
		}

	}

	longest := ""

	for type_name := range types {

		if strings.HasPrefix(name, type_name+"From") == true {

			if len(type_name) > len(longest) {
				longest = type_name
			}

		}

	}

	if longest != "" {
		return longest
	}

	if strings.HasPrefix(name, "New") == true || strings.HasPrefix(name, "Parse") == true || strings.HasPrefix(name, "new") == true {
		return declaredResultType(func_decl, types)
	}

	return ""

}

func freeFunctions(file *File) []*ast.FuncDecl {

	result := make([]*ast.FuncDecl, 0)

	for _, decl := range file.AST.Decls {

		func_decl, ok := decl.(*ast.FuncDecl)

		if ok == true && func_decl.Recv == nil && func_decl.Name.Name != "init" {
			result = append(result, func_decl)
		}

	}

	return result

}

func declaredTypes(file *File) []*ast.TypeSpec {

	result := make([]*ast.TypeSpec, 0)

	for _, decl := range file.AST.Decls {

		gen_decl, ok := decl.(*ast.GenDecl)

		if ok == true && gen_decl.Tok == token.TYPE {

			for _, spec := range gen_decl.Specs {

				type_spec, ok := spec.(*ast.TypeSpec)

				if ok == true {
					result = append(result, type_spec)
				}

			}

		}

	}

	return result

}

func checkOneType(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	}

	names := make([]string, 0)

	for _, type_spec := range declaredTypes(file) {

		if isExported(type_spec.Name.Name) == true {
			names = append(names, type_spec.Name.Name)
		}

	}

	result := make([]Diagnostic, 0)

	for index, name := range names {

		if index > 0 {
			result = append(
				result,
				file.Diagnostic(
					"layout/one-type",
					file.AST.Pos(),
					"file declares %d exported types (%s); %s belongs in %s.go",
					len(names),
					strings.Join(names, ", "),
					name,
					name,
				),
			)
		}

	}

	return result

}

func checkTypeFile(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	}

	types := declaredTypes(file)

	if len(types) == 0 {
		return nil
	}

	has_exported := false

	for _, type_spec := range types {

		if isExported(type_spec.Name.Name) == true {
			has_exported = true
		}

	}

	if has_exported == true {
		return nil
	}

	base := identifierBaseName(file)
	result := make([]Diagnostic, 0)

	for _, type_spec := range types {

		if type_spec.Name.Name != base {
			result = append(
				result,
				file.Diagnostic(
					"layout/type-file",
					type_spec.Pos(),
					"unexported type %s must live in %s.go",
					type_spec.Name.Name,
					type_spec.Name.Name,
				),
			)
		}

	}

	return result

}

func checkOneFunction(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true || file.AST.Name.Name == "main" {
		return nil
	}

	if oneFunctionExempt[file.Rel] == true {
		return nil
	}

	candidates := make([]*ast.FuncDecl, 0)

	for _, func_decl := range freeFunctions(file) {

		if ownedTypeName(func_decl, file.Pkg.Types) == "" {
			candidates = append(candidates, func_decl)
		}

	}

	if len(candidates) <= 1 {
		return nil
	}

	base := identifierBaseName(file)
	matching := 0

	for _, func_decl := range candidates {

		if func_decl.Name.Name == base {
			matching++
		}

	}

	result := make([]Diagnostic, 0)

	for _, func_decl := range candidates {

		if matching == 1 && func_decl.Name.Name == base {
			continue
		}

		result = append(
			result,
			file.Diagnostic(
				"layout/one-function",
				func_decl.Pos(),
				"file declares %d free functions; move %s into %s.go",
				len(candidates),
				func_decl.Name.Name,
				func_decl.Name.Name,
			),
		)

	}

	return result

}

func checkFunctionFile(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true || file.AST.Name.Name == "main" {
		return nil
	}

	candidates := make([]*ast.FuncDecl, 0)

	for _, func_decl := range freeFunctions(file) {

		if ownedTypeName(func_decl, file.Pkg.Types) == "" {
			candidates = append(candidates, func_decl)
		}

	}

	if len(candidates) != 1 {
		return nil
	}

	base := identifierBaseName(file)

	if candidates[0].Name.Name == base {
		return nil
	}

	return []Diagnostic{
		file.Diagnostic(
			"layout/function-file",
			candidates[0].Pos(),
			"%s must live in %s.go",
			candidates[0].Name.Name,
			candidates[0].Name.Name,
		),
	}

}

func checkInit(_ *Context, file *File) []Diagnostic {

	count := 0

	for _, decl := range file.AST.Decls {

		func_decl, ok := decl.(*ast.FuncDecl)

		if ok == true && func_decl.Recv == nil && func_decl.Name.Name == "init" {
			count++
		}

	}

	if count <= 1 {
		return nil
	}

	return []Diagnostic{
		file.Diagnostic("layout/init", file.AST.Pos(), "file declares %d init() functions", count),
	}

}

func checkToolsRoot(_ *Context, file *File) []Diagnostic {

	if path.Dir(file.Rel) != "tools" {
		return nil
	}

	if file.Rel != "tools/Toolset.go" && file.Rel != "tools/Toolset_test.go" {
		return []Diagnostic{
			file.Diagnostic(
				"layout/tools-root",
				file.AST.Pos(),
				"source/tools may only contain Toolset.go and Toolset_test.go",
			),
		}
	}

	return nil

}

func checkUtilsMirror(context *Context, file *File) []Diagnostic {

	if strings.HasPrefix(file.Rel, "utils/") == false {
		return nil
	}

	parts := strings.Split(file.Rel, "/")

	if len(parts) < 3 {
		return nil
	}

	name := parts[1]

	for _, pkg := range context.Packages {

		if pkg.Rel == name && pkg.Name == name {
			return []Diagnostic{
				file.Diagnostic(
					"layout/utils-mirror",
					file.AST.Pos(),
					"utils/%s duplicates the top-level package %s",
					name,
					name,
				),
			}
		}

	}

	return nil

}

func checkMethodFile(_ *Context, pkg *Package) []Diagnostic {

	result := make([]Diagnostic, 0)

	for _, file := range pkg.Files {

		for _, decl := range file.AST.Decls {

			func_decl, ok := decl.(*ast.FuncDecl)

			if ok == false || func_decl.Recv == nil || len(func_decl.Recv.List) == 0 {
				continue
			}

			type_name := receiverTypeName(func_decl.Recv.List[0].Type)

			if type_name == "" {
				continue
			}

			owner := pkg.Types[type_name]

			if owner == nil || owner == file {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"layout/method-file",
					func_decl.Pos(),
					"method (%s *%s).%s must live in %s.go",
					receiverName(func_decl.Recv.List[0]),
					type_name,
					func_decl.Name.Name,
					type_name,
				),
			)

		}

	}

	return result

}

func checkConstructorFile(_ *Context, pkg *Package) []Diagnostic {

	result := make([]Diagnostic, 0)

	for _, file := range pkg.Files {

		if file.IsTest == true {
			continue
		}

		for _, func_decl := range freeFunctions(file) {

			type_name := ownedTypeName(func_decl, pkg.Types)

			if type_name == "" {
				continue
			}

			if identifierBaseName(file) == type_name {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"layout/constructor-file",
					func_decl.Pos(),
					"%s builds %s and must live in %s.go",
					func_decl.Name.Name,
					type_name,
					type_name,
				),
			)

		}

	}

	return result

}

func hasPlatformSuffix(base string) bool {

	parts := strings.Split(base, "_")

	if len(parts) < 2 {
		return false
	}

	known := map[string]bool{
		"aix": true, "android": true, "darwin": true, "dragonfly": true,
		"freebsd": true, "illumos": true, "ios": true, "js": true,
		"linux": true, "netbsd": true, "openbsd": true, "other": true,
		"plan9": true, "solaris": true, "unix": true, "wasip1": true,
		"windows": true,
		"386":     true, "amd64": true, "arm": true, "arm64": true,
		"loong64": true, "mips": true, "mips64": true, "mips64le": true,
		"mipsle": true, "ppc64": true, "ppc64le": true, "riscv64": true,
		"s390x": true, "wasm": true,
	}

	for index, part := range parts {

		if index == 0 {
			continue
		}

		if known[part] == false {
			return false
		}

	}

	return true

}
