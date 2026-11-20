package requirements

import "exocomp/types"
import utils_ast "exocomp/utils/ast"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "strings"

func (tool *Requirements) DefineStruct(path string, symbol string, declaration string, behavior string) (string, error) {

	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

		if err2 == nil {

			declaration = strings.TrimSpace(declaration)

			if strings.HasPrefix(declaration, "struct{") || strings.HasPrefix(declaration, "struct {") {
				declaration = "type " + symbol + " struct " + strings.TrimSpace(declaration[6:])
			}

			prettified := utils_ast.GetSymbol([]byte(strings.Join([]string{
				"package dummy",
				declaration,
			}, "\n")), symbol, "struct")

			if prettified != nil {

				readRequirements(tool)

				tool.mutex.Lock()

				_, ok3 := tool.contents[internal_path]

				if ok3 == false {
					tool.contents[internal_path] = make(map[string]types.Requirement)
				}

				tool.contents[internal_path][symbol] = types.Requirement{
					File:        internal_path,
					Type:        "struct",
					Declaration: prettified.Body,
					Symbol:      prettified.Name,
					Behavior:    behavior,
				}

				tool.mutex.Unlock()

				err4 := writeRequirements(tool)

				if err4 == nil {
					return fmt.Sprintf("requirements.DefineStruct: Symbol \"%s\" defined as \"%s\"", prettified.Name, prettified.Body), nil
				} else {
					return "", fmt.Errorf("requirements.DefineStruct: %s", err4.Error())
				}

			} else {
				return "", fmt.Errorf("requirements.DefineStruct: Invalid Go syntax. \"type %s struct { ...}\" must be defined!", symbol)
			}

		} else {
			return "", fmt.Errorf("requirements.DefineStruct: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.DefineStruct: %s", err1.Error())
	}

}
