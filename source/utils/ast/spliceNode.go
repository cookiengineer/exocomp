package ast

import "go/ast"
import "go/token"

func spliceNode(source []byte, fileset *token.FileSet, node ast.Node, replacement string) []byte {

	start := fileset.Position(node.Pos()).Offset
	end := fileset.Position(node.End()).Offset

	result := make([]byte, 0, len(source)+len(replacement))
	result = append(result, source[0:start]...)
	result = append(result, []byte(replacement)...)
	result = append(result, source[end:]...)

	return result

}

