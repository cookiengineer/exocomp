package lint

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
