package requirements

import "exocomp/types"
import utils_ast "exocomp/utils/ast"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "strings"

func (tool *Requirements) DefineFunc(path string, symbol string, declaration string, behavior string) (string, error) {

	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

		if err2 == nil {

			declaration = strings.TrimSpace(declaration)

			prettified := utils_ast.GetSymbol([]byte(strings.Join([]string{
				"package dummy",
				declaration,
			}, "\n")), symbol, "func")

			if prettified != nil {

				readRequirements(tool)

				tool.mutex.Lock()

				_, ok1 := tool.contents[internal_path]

				if ok1 == false {
					tool.contents[internal_path] = make(map[string]types.Requirement)
				}

				tool.contents[internal_path][symbol] = types.Requirement{
					File:        internal_path,
					Type:        "func",
					Declaration: prettified.Body,
					Symbol:      prettified.Name,
					Behavior:    behavior,
				}

				tool.mutex.Unlock()

				err4 := writeRequirements(tool)

				if err4 == nil {
					return fmt.Sprintf("requirements.DefineFunc: Symbol \"%s\" defined as \"%s\"", prettified.Name, prettified.Body), nil
				} else {
					return "", fmt.Errorf("requirements.DefineFunc: %s", err4.Error())
				}

			} else {
				return "", fmt.Errorf("requirements.DefineFunc: Invalid Go syntax. \"func %s (...) (...)\" must be defined!", symbol)
			}

		} else {
			return "", fmt.Errorf("requirements.DefineFunc: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.DefineFunc: %s", err1.Error())
	}

}
