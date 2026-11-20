package lint

import "path/filepath"
import "strings"

func checkFileName(_ *Context, file *File) []Diagnostic {

	if file.IsTest == true {
		return nil
	} else {

		base := fileBaseName(file)

		if strings.Contains(base, "_") == false {
			return nil
		} else if hasPlatformSuffix(base) == true {
			return nil
		} else {
			return []Diagnostic{
				file.Diagnostic(
					"naming/file",
					file.AST.Pos(),
					"file name %q uses underscores; name it after the identifier it contains",
					filepath.Base(file.Rel),
				),
			}
		}

	}

}
