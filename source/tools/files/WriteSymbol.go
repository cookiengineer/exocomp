package files

import utils_ast "exocomp/utils/ast"
import "bytes"
import "fmt"
import "os"
import "strings"

func (tool *Files) WriteSymbol(path string, symbol string, declaration string) (string, error) {

	resolved, err0 := resolveFilesPath(tool.Sandbox, tool.Capabilities, path, true)

	if err0 == nil {

		source, err1 := os.ReadFile(resolved)

		if err1 == nil {

			declaration_type := ""

			if strings.HasPrefix(strings.TrimSpace(declaration), "func") {
				declaration_type = "func"
			}

			result := utils_ast.WriteSymbol(source, symbol, declaration, declaration_type)

			if bytes.Equal(result, source) == true {
				return "", fmt.Errorf("files.WriteSymbol: File \"%s\" has no Symbol \"%s\"", path, symbol)
			}

			err2 := os.WriteFile(resolved, result, 0666)

			if err2 == nil {
				return fmt.Sprintf("files.WriteSymbol: File \"%s\" Symbol \"%s\" written.", path, symbol), nil
			} else {
				return sanitizeFilesystemError("files", "WriteSymbol", "File", path, err2)
			}

		} else {
			return sanitizeFilesystemError("files", "WriteSymbol", "File", path, err1)
		}

	} else {
		return "", fmt.Errorf("files.WriteSymbol: %s", err0.Error())
	}

}
