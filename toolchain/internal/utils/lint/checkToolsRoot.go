package lint

import "path"

func checkToolsRoot(_ *Context, file *File) []Diagnostic {

	if path.Dir(file.Rel) != "tools" {
		return nil
	} else if file.Rel != "tools/Toolset.go" && file.Rel != "tools/Toolset_test.go" {
		return []Diagnostic{
			file.Diagnostic(
				"layout/tools-root",
				file.AST.Pos(),
				"source/tools may only contain Toolset.go and Toolset_test.go",
			),
		}
	} else {
		return nil
	}

}
