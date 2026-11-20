package bugs

import utils_sandbox "exocomp/utils/sandbox"
import "fmt"

func (tool *Bugs) Fix(path string, symbol string) (string, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	tmp1, err1 := utils_sandbox.ResolvePath(tool.Sandbox, path)

	if err1 == nil {

		internal_path, err2 := utils_sandbox.SanitizePath(tool.Playground, tmp1)

		if err2 == nil {

			readBugs(tool)

			_, ok1 := tool.contents[internal_path]

			if ok1 == true {

				bug_report, ok2 := tool.contents[internal_path][symbol]

				if ok2 == true {

					bug_report.IsFixed = true
					tool.contents[internal_path][symbol] = bug_report

					err3 := writeBugs(tool)

					if err3 == nil {
						return fmt.Sprintf("bugs.Fix: Bug report marked as fixed."), nil
					} else {
						return "", fmt.Errorf("bugs.Fix: %s", err3.Error())
					}

				} else {
					return "", fmt.Errorf("bugs.Fix: No bug report available for path \"%s\" and symbol \"%s\"", path, symbol)
				}

			} else {
				return "", fmt.Errorf("bugs.Fix: No bug reports available for path \"%s\".", path)
			}

		} else {
			return "", fmt.Errorf("bugs.Fix: %s", err2.Error())
		}

	} else {
		return "", fmt.Errorf("bugs.Fix: %s", err1.Error())
	}

}
