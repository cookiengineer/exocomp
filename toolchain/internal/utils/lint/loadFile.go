package lint

import "go/parser"
import "os"
import "path/filepath"
import "strings"

func loadFile(context *Context, root string, path string) (*File, error) {

	source, err0 := os.ReadFile(path)

	if err0 == nil {

		parsed, err1 := parser.ParseFile(context.Fset, path, source, parser.ParseComments|parser.AllErrors)

		if err1 == nil {

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

		} else {
			return nil, err1
		}

	} else {
		return nil, err0
	}

}
