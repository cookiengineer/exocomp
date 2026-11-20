package requirements

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "sort"
import "strings"

func (tool *Requirements) List() (string, error) {

	readRequirements(tool)

	lines := make([]string, 0)

	tool.mutex.RLock()

	for _, specifications := range tool.contents {

		for _, specification := range specifications {

			resolved_path, err1 := utils_sandbox.ResolvePath(tool.Playground, specification.File)

			if err1 == nil {

				sandbox_path, err2 := utils_sandbox.SanitizePath(tool.Sandbox, resolved_path)

				if err2 == nil {
					lines = append(lines, fmt.Sprintf("- File: \"%s\", Symbol: \"%s\", Declaration: \"%s\", Behavior: \"%s\", IsImplemented: %t", sandbox_path, specification.Symbol, specification.Declaration, specification.Behavior, specification.IsImplemented))
				}

			}

		}

	}

	tool.mutex.RUnlock()

	sort.Strings(lines)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("requirements.List: %d specifications.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}
