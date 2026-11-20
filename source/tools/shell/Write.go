package shell

import utils_shell "exocomp/utils/shell"
import "fmt"

func (tool *Shell) Write(session_id string, data string) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		if len(data) <= shellMaxWrite {

			written, err0 := session.write([]byte(data))

			if err0 == nil {

				status, exit_code := session.status()

				return fmt.Sprintf("shell.Write: Session \"%s\" wrote %d bytes, status \"%s\"%s.", session_id, written, status, utils_shell.FormatExitCode(status, exit_code)), nil

			} else {
				return "", fmt.Errorf("shell.Write: Cannot write to session \"%s\".", session_id)
			}

		} else {
			return "", fmt.Errorf("shell.Write: Invalid parameter \"data\" exceeds the %d byte limit.", shellMaxWrite)
		}

	} else {
		return "", fmt.Errorf("shell.Write: Session \"%s\" not found.", session_id)
	}

}
