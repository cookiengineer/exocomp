package requirements

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "sort"
import "strings"

func (tool *Requirements) Search(path string, symbol string) (string, error) {

	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		search_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

		if err2 == nil {

			readRequirements(tool)

			lines := make([]string, 0)

			tool.mutex.RLock()

			for internal_path, specifications := range tool.contents {

				if strings.HasPrefix(internal_path, search_path) {

					for internal_symbol, specification := range specifications {

						if symbol != "" {

							matches_symbol := strings.Contains(strings.ToLower(internal_symbol), strings.ToLower(symbol))

							if matches_symbol == true {

								sandbox_path, err3 := utils_sandbox.SanitizePath(tool.Sandbox, specification.File)

								if err3 == nil {
									lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Declaration: %s, Behavior: %s", sandbox_path, specification.Symbol, specification.Declaration, specification.Behavior))
								}

							}

						} else if symbol == "" {

							sandbox_path, err3 := utils_sandbox.SanitizePath(tool.Sandbox, specification.File)

							if err3 == nil {
								lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Declaration: %s, Behavior: %s", sandbox_path, specification.Symbol, specification.Declaration, specification.Behavior))
							}

						}

					}

				}

			}

			tool.mutex.RUnlock()

			sort.Strings(lines)

			result := make([]string, 0)

			if symbol != "" {
				result = append(result, fmt.Sprintf("requirements.Search: %s#%s* contains %d specifications.", path, symbol, len(lines)))
			} else {
				result = append(result, fmt.Sprintf("requirements.Search: %s#* contains %d specifications.", path, len(lines)))
			}

			for l := 0; l < len(lines); l++ {
				result = append(result, lines[l])
			}

			return strings.Join(result, "\n"), nil

		} else {
			return "", fmt.Errorf("requirements.Search: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("requirements.Search: %s", err1.Error())
	}

}
