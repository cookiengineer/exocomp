package shell

import utils_shell "exocomp/utils/shell"
import "fmt"

func (tool *Shell) Close(session_id string, force bool) (string, error) {

	tool.mutex.Lock()
	session, ok := tool.sessions[session_id]

	if ok == true {
		delete(tool.sessions, session_id)
	}

	tool.mutex.Unlock()

	if ok == true && session != nil {

		data, status, exit_code := session.close(force)
		total, _, _ := session.outputStats()

		header := fmt.Sprintf("shell.Close: Session \"%s\" closed with status \"%s\"%s, log %d bytes.", session_id, status, utils_shell.FormatExitCode(status, exit_code), total)

		if len(data) > 0 {
			return header + "\n" + string(data), nil
		} else {
			return header + "\n<no new output>", nil
		}

	} else {
		return "", fmt.Errorf("shell.Close: Session \"%s\" not found.", session_id)
	}

}
