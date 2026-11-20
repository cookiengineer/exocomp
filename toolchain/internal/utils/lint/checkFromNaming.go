package lint

import "strings"

func checkFromNaming(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		result := make([]Diagnostic, 0)

		for _, func_decl := range freeFunctions(file) {

			name := func_decl.Name.Name

			if strings.HasPrefix(name, "From") == false || len(name) <= len("From") {
				continue
			}

			next := name[len("From")]

			if next < 'A' || next > 'Z' {
				continue
			}

			result = append(
				result,
				file.Diagnostic(
					"naming/from",
					func_decl.Pos(),
					"converter %s must be named <Type>From<Source>",
					name,
				),
			)

		}

		return result

	}

}
