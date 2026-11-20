package shell

import utils_shell "exocomp/utils/shell"
import "fmt"

func (tool *Shell) callOpen(program string, arguments []string, workdir string, raw map[string]interface{}) (string, error) {

	columns, err0 := utils_shell.ParseInt(raw, "columns", shellDefaultColumns)

	if err0 == nil {

		rows, err1 := utils_shell.ParseInt(raw, "rows", shellDefaultRows)

		if err1 == nil {
			return tool.Open(program, arguments, workdir, columns, rows)
		} else {
			return "", fmt.Errorf("shell.Open: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("shell.Open: %s", err0.Error())
	}

}
