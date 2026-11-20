package main

import "fmt"
import "go/ast"
import "go/token"
import "strconv"
import "strings"

// FileRule checks a single parsed file.
type FileRule struct {
	ID      string
	Summary string
	Check   func(*Context, *File) []Diagnostic
}

// PackageRule checks a whole directory.
type PackageRule struct {
	ID      string
	Summary string
	Check   func(*Context, *Package) []Diagnostic
}

var fileRules []FileRule
var packageRules []PackageRule

func init() {

	fileRules = []FileRule{
		{"dep/tool-isolation", "tool packages must not import each other", checkToolIsolation},
		{"dep/tools-allowed", "tool packages may only import types, schemas, agents and utils", checkToolsAllowed},
		{"dep/types-no-capabilities", "types must not import utils/capabilities", checkTypesNoCapabilities},
		{"dep/utils-no-tools", "utils must not import tools", checkUtilsNoTools},
		{"branch/guard", "terminal guard clauses become positive-first branches", checkGuard},
		{"branch/guard-chain", "consecutive guards become one else-if chain", checkGuardChain},
		{"cond/bare", "conditions are compared explicitly", checkBareCondition},
		{"cond/not", "no ! negation in conditions", checkNot},
		{"errors/lowercase", "error messages start with a capital letter", checkErrorLowercase},
		{"errors/prefix", "error messages carry a component prefix", checkErrorPrefix},
		{"import/group", "project imports first, then stdlib, then third-party", checkImportGroup},
		{"import/order", "imports are sorted alphabetically within their group", checkImportOrder},
		{"import/paren", "one import per line, no parenthesized block", checkImportParen},
		{"import/testing-last", "testing is the last stdlib import in tests", checkTestingLast},
		{"layout/function-file", "a lone free function lives in a file named after it", checkFunctionFile},
		{"layout/init", "at most one init() per file", checkInit},
		{"layout/one-function", "every free function lives in its own file", checkOneFunction},
		{"layout/one-type", "one primary type per file", checkOneType},
		{"layout/tools-root", "source/tools contains only Toolset.go", checkToolsRoot},
		{"layout/type-file", "a lone unexported type lives in a file named after it", checkTypeFile},
		{"layout/utils-mirror", "utils/<name> must not mirror a top-level package", checkUtilsMirror},
		{"naming/boolean", "boolean values read as predicates", checkBooleanNaming},
		{"naming/file", "Go file names use the identifier spelling, no underscores", checkFileName},
		{"naming/from", "converters are named <Type>From<Source>", checkFromNaming},
		{"naming/func", "unexported functions use camelCase, never snake_case", checkFuncNaming},
		{"naming/type-initialism", "type names are spelled out, not collapsed to initials", checkTypeInitialism},
		{"style/any", "use any instead of interface{}", checkAny},
		{"style/blank-block", "blank lines around blocks and above switch cases", checkBlankBlock},
		{"style/gofmt", "file is formatted", checkGofmt},
		{"style/json-indent", "JSON indentation uses tabs", checkJSONIndent},
		{"style/local-case", "local variables use snake_case", checkLocalCase},
		{"test/file", "test files are named after a source file", checkTestFile},
		{"test/func", "test functions are named Test<Type>_<Method>", checkTestFunc},
	}

	packageRules = []PackageRule{
		{"layout/constructor-file", "New<Type> and Parse<Type> live in the type file", checkConstructorFile},
		{"layout/method-file", "methods live in the file of their receiver type", checkMethodFile},
		{"naming/receiver", "one receiver name per type", checkReceiverNaming},
	}

}

// Position returns the token position of pos.
func (file *File) Position(pos token.Pos) token.Position {
	return file.Fset.Position(pos)
}

// Line returns the 1-based line of pos.
func (file *File) Line(pos token.Pos) int {
	return file.Fset.Position(pos).Line
}

// Text returns the source slice covered by node.
func (file *File) Text(node ast.Node) string {

	start := file.Fset.Position(node.Pos()).Offset
	end := file.Fset.Position(node.End()).Offset

	if start < 0 || end > len(file.Src) || start > end {
		return ""
	}

	return strings.TrimSpace(string(file.Src[start:end]))

}

// Diagnostic builds a diagnostic for rule at pos.
func (file *File) Diagnostic(rule string, pos token.Pos, format string, arguments ...any) Diagnostic {
	return Diagnostic{
		Position: file.Position(pos),
		Rule:     rule,
		Message:  fmt.Sprintf(format, arguments...),
	}
}

// Package.Diagnostic builds a diagnostic anchored to a file position.
func (pkg *Package) Diagnostic(rule string, file *File, pos token.Pos, format string, arguments ...any) Diagnostic {
	return file.Diagnostic(rule, pos, format, arguments...)
}

func isExported(name string) bool {

	if name == "" {
		return false
	}

	return name[0] >= 'A' && name[0] <= 'Z'

}

func receiverTypeName(expr ast.Expr) string {

	switch typed := expr.(type) {
	case *ast.StarExpr:
		return receiverTypeName(typed.X)
	case *ast.Ident:
		return typed.Name
	case *ast.IndexExpr:
		return receiverTypeName(typed.X)
	case *ast.IndexListExpr:
		return receiverTypeName(typed.X)
	default:
		return ""
	}

}

func importPath(spec *ast.ImportSpec) string {

	if spec == nil || spec.Path == nil {
		return ""
	}

	value, err := strconv.Unquote(spec.Path.Value)

	if err == nil {
		return value
	} else {
		return ""
	}

}

func unwrapParens(expr ast.Expr) ast.Expr {

	for {

		typed, ok := expr.(*ast.ParenExpr)

		if ok == true {
			expr = typed.X
		} else {
			return expr
		}

	}

}

func receiverName(field *ast.Field) string {

	if len(field.Names) == 0 {
		return ""
	}

	return field.Names[0].Name

}
