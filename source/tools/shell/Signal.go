package shell

import utils_shell "exocomp/utils/shell"
import "fmt"

func (tool *Shell) Signal(session_id string, name string) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		signal, err0 := session.signal(name)

		if err0 == nil {

			status, exit_code := session.status()

			return fmt.Sprintf("shell.Signal: Session \"%s\" sent SIG%s, status \"%s\"%s.", session_id, signal, status, utils_shell.FormatExitCode(status, exit_code)), nil

		} else {
			return "", fmt.Errorf("shell.Signal: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("shell.Signal: Session \"%s\" not found.", session_id)
	}

}
