package lint

import "go/format"

func checkGofmt(_ *Context, file *File) []Diagnostic {

	formatted, err := format.Source(file.Src)

	if err != nil {
		return []Diagnostic{
			file.Diagnostic("style/gofmt", file.AST.Pos(), "cannot format: %s", err.Error()),
		}
	} else if string(formatted) == string(file.Src) {
		return nil
	} else {
		return []Diagnostic{
			file.Diagnostic("style/gofmt", file.AST.Pos(), "file is not gofmt formatted"),
		}
	}

}
