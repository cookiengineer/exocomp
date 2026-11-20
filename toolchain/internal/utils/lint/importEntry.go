package lint

import "go/ast"

type importEntry struct {
	Path string
	Spec *ast.ImportSpec
	Line int
}
