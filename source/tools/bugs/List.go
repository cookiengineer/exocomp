package bugs

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"
import "sort"
import "strings"

func (tool *Bugs) List() (string, error) {

	readBugs(tool)

	lines := make([]string, 0)

	tool.mutex.RLock()

	for _, bug_reports := range tool.contents {

		for _, bug_report := range bug_reports {

			if bug_report.IsFixed == false {

				resolved_path, err1 := utils_sandbox.ResolvePath(tool.Playground, bug_report.File)

				if err1 == nil {

					sandbox_path, err2 := utils_sandbox.SanitizePath(tool.Sandbox, resolved_path)

					if err2 == nil {
						lines = append(lines, fmt.Sprintf("- File: %s, Symbol: %s, Description: %s", sandbox_path, bug_report.Symbol, bug_report.Description))
					}

				}

			}

		}

	}

	tool.mutex.RUnlock()

	sort.Strings(lines)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("bugs.List: %d unfixed bug reports.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}
