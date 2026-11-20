package requirements

import utils_ast "exocomp/utils/ast"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "os"

func (tool *Requirements) Signoff(path string, symbol string) (string, error) {

	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

		if err2 == nil {

			readRequirements(tool)

			tool.mutex.RLock()
			specification, ok1 := tool.contents[internal_path][symbol]
			tool.mutex.RUnlock()

			if ok1 == true {

				resolved_path, err3 := utils_sandbox.ResolvePath(tool.Playground, specification.File)

				if err3 == nil {

					source, err4 := os.ReadFile(resolved_path)

					if err4 == nil {

						if utils_ast.HasSymbol(source, specification.Symbol, specification.Type) == true {

							tool.mutex.Lock()
							specification.IsImplemented = true
							tool.contents[internal_path][symbol] = specification
							tool.mutex.Unlock()

							err5 := writeRequirements(tool)

							if err5 == nil {
								return fmt.Sprintf("requirements.Signoff: %s#%s marked as implemented.", path, symbol), nil
							} else {
								return "", fmt.Errorf("requirements.Signoff: %s", err5.Error())
							}

						} else {
							return "", fmt.Errorf("requirements.Signoff: Symbol \"%s\" is not implemented in \"%s\" yet.", symbol, path)
						}

					} else {
						return "", fmt.Errorf("requirements.Signoff: File \"%s\" does not exist.", path)
					}

				} else {
					return "", fmt.Errorf("requirements.Signoff: %s", err3.Error())
				}

			} else {
				return "", fmt.Errorf("requirements.Signoff: No specification available for path \"%s\" and symbol \"%s\"", path, symbol)
			}

		} else {
			return "", fmt.Errorf("requirements.Signoff: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.Signoff: %s", err1.Error())
	}

}
