package bugs

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "sort"
import "strings"

func (tool *Bugs) Search(path string, symbol string) (string, error) {

	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

		if err2 == nil {

			readBugs(tool)

			if symbol != "" {

				lines := make([]string, 0)

				tool.mutex.RLock()

				_, ok1 := tool.contents[internal_path]

				if ok1 == true {

					bug_report, ok2 := tool.contents[internal_path][symbol]

					if ok2 == true {

						if bug_report.IsFixed == false {

							sandbox_path, err3 := utils_sandbox.SanitizePath(tool.Sandbox, bug_report.File)

							if err3 == nil {
								lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Description: %s", sandbox_path, bug_report.Symbol, bug_report.Description))
							}

						}

					}

				}

				tool.mutex.RUnlock()

				sort.Strings(lines)

				result := make([]string, 0)
				result = append(result, fmt.Sprintf("bugs.Search: %s#%s contains %d bug reports.", path, symbol, len(lines)))

				for l := 0; l < len(lines); l++ {
					result = append(result, lines[l])
				}

				return strings.Join(result, "\n"), nil

			} else {

				lines := make([]string, 0)

				tool.mutex.RLock()

				bug_reports, ok1 := tool.contents[internal_path]

				if ok1 == true {

					for _, bug_report := range bug_reports {

						if bug_report.IsFixed == false {

							sandbox_path, err3 := utils_sandbox.SanitizePath(tool.Sandbox, bug_report.File)

							if err3 == nil {
								lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Description: %s", sandbox_path, bug_report.Symbol, bug_report.Description))
							}

						}

					}

				}

				tool.mutex.RUnlock()

				sort.Strings(lines)

				result := make([]string, 0)
				result = append(result, fmt.Sprintf("bugs.Search: %s contains %d bug reports.", path, len(lines)))

				for l := 0; l < len(lines); l++ {
					result = append(result, lines[l])
				}

				return strings.Join(result, "\n"), nil

			}

		} else {
			return "", fmt.Errorf("bugs.Search: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("bugs.Search: %s", err1.Error())
	}

}
