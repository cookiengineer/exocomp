package lint

func importPaths(file *File) []string {

	result := make([]string, 0, len(file.AST.Imports))

	for _, spec := range file.AST.Imports {

		value := importPath(spec)

		if value != "" {
			result = append(result, value)
		}

	}

	return result

}
