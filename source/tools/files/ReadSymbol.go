package files

import utils_ast "exocomp/utils/ast"
import "fmt"
import "os"
import "strings"

func (tool *Files) ReadSymbol(path string, symbol string) (string, error) {

	resolved, err0 := resolveFilesPath(tool.Sandbox, tool.Capabilities, path, false)

	if err0 == nil {

		bytes, err1 := os.ReadFile(resolved)

		if err1 == nil {

			found := utils_ast.GetSymbol(bytes, symbol, "")

			if found != nil {

				result := strings.Join([]string{
					fmt.Sprintf("files.ReadSymbol: File \"%s\" Symbol \"%s\" with Type \"%s\"", path, found.Name, found.Type),
					fmt.Sprintf("%s", found.Body),
				}, "\n")

				return result, nil

			} else {
				return "", fmt.Errorf("files.ReadSymbol: File \"%s\" has no Symbol \"%s\"", path, symbol)
			}

		} else {
			return sanitizeFilesystemError("files", "ReadSymbol", "File", path, err1)
		}

	} else {
		return "", fmt.Errorf("files.ReadSymbol: %s", err0.Error())
	}

}
