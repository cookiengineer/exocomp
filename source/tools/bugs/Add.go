package bugs

import "exocomp/types"
import utils_sandbox "exocomp/utils/sandbox"
import "fmt"

func (tool *Bugs) Add(path string, symbol string, description string) (string, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

		if err2 == nil {

			readBugs(tool)

			_, ok1 := tool.contents[internal_path]

			if ok1 == false {
				tool.contents[internal_path] = make(map[string]types.Bug)
			}

			_, ok2 := tool.contents[internal_path][symbol]

			if ok2 == false {

				tool.contents[internal_path][symbol] = types.Bug{
					IsFixed:     false,
					File:        internal_path,
					Symbol:      symbol,
					Description: description,
				}

				err3 := writeBugs(tool)

				if err3 == nil {
					return fmt.Sprintf("bugs.Add: Bug report with %d B written.", len(description)), nil
				} else {
					return "", fmt.Errorf("bugs.Add: %s", err3.Error())
				}

			} else {

				bug_report := tool.contents[internal_path][symbol]
				bug_report.Description = description
				bug_report.IsFixed = false
				tool.contents[internal_path][symbol] = bug_report

				err3 := writeBugs(tool)

				if err3 == nil {
					return fmt.Sprintf("bugs.Add: Bug report with %d B updated.", len(description)), nil
				} else {
					return "", fmt.Errorf("bugs.Add: %s", err3.Error())
				}

			}

		} else {
			return "", fmt.Errorf("bugs.Add: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("bugs.Add: %s", err1.Error())
	}

}
