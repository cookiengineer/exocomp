package requirements

import "exocomp/types"
import utils_ast "exocomp/utils/ast"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "strings"

func (tool *Requirements) DefineInterface(path string, symbol string, declaration string, behavior string) (string, error) {

	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

		if err2 == nil {

			declaration = strings.TrimSpace(declaration)

			if strings.HasPrefix(declaration, "interface{") || strings.HasPrefix(declaration, "interface {") {
				declaration = "type " + symbol + " interface " + strings.TrimSpace(declaration[9:])
			}

			prettified := utils_ast.GetSymbol([]byte(strings.Join([]string{
				"package dummy",
				declaration,
			}, "\n")), symbol, "interface")

			if prettified != nil {

				readRequirements(tool)

				tool.mutex.Lock()

				_, ok3 := tool.contents[internal_path]

				if ok3 == false {
					tool.contents[internal_path] = make(map[string]types.Requirement)
				}

				tool.contents[internal_path][symbol] = types.Requirement{
					File:        internal_path,
					Type:        "interface",
					Declaration: prettified.Body,
					Symbol:      prettified.Name,
					Behavior:    behavior,
				}

				tool.mutex.Unlock()

				err4 := writeRequirements(tool)

				if err4 == nil {
					return fmt.Sprintf("requirements.DefineInterface: Symbol \"%s\" defined as \"%s\"", prettified.Name, prettified.Body), nil
				} else {
					return "", fmt.Errorf("requirements.DefineInterface: %s", err4.Error())
				}

			} else {
				return "", fmt.Errorf("requirements.DefineInterface: Invalid Go syntax. \"type %s interface { ...}\" must be defined!", symbol)
			}

		} else {
			return "", fmt.Errorf("requirements.DefineInterface: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.DefineInterface: %s", err1.Error())
	}

}
