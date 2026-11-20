package shell

import utils_shell "exocomp/utils/shell"
import "fmt"
import "time"

func (tool *Shell) Read(session_id string, wait_ms int, limit int) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		if wait_ms >= 0 && wait_ms <= 30000 {

			if limit >= 1 && limit <= shellMaxRead {

				data, status, exit_code := session.read(time.Duration(wait_ms)*time.Millisecond, limit)
				dropped, total, truncated := session.outputStats()
				pending := session.pendingOutput()

				header := fmt.Sprintf("shell.Read: Session \"%s\" status \"%s\" read %d bytes (total %d, dropped %d)%s.", session_id, status, len(data), total, dropped, utils_shell.FormatExitCode(status, exit_code))

				if truncated == true || pending > 0 || dropped > 0 {
					header += " Output was truncated; use shell.Log."
				}

				if len(data) > 0 {
					return header + "\n" + string(data), nil
				} else {
					return header + "\n<no new output>", nil
				}

			} else {
				return "", fmt.Errorf("shell.Read: Invalid parameter \"limit\" must be between 1 and %d.", shellMaxRead)
			}

		} else {
			return "", fmt.Errorf("shell.Read: Invalid parameter \"wait_ms\" must be between 0 and 30000.")
		}

	} else {
		return "", fmt.Errorf("shell.Read: Session \"%s\" not found.", session_id)
	}

}
