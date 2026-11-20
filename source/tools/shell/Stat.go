package shell

import utils_shell "exocomp/utils/shell"
import "fmt"
import "strings"
import "time"

func (tool *Shell) Stat(session_id string) (string, error) {

	session := tool.getSession(session_id)

	if session != nil {

		session.mutex.Lock()
		program := session.Program
		arguments := append([]string{}, session.Arguments...)
		workdir := session.Workdir
		status := session.Status
		exit_code := session.ExitCode
		started_at := session.StartedAt
		finished_at := session.FinishedAt
		last_activity := session.LastActivity
		columns := session.Columns
		rows := session.Rows
		log_size := session.logSize
		dropped := session.dropped
		log_truncated := session.logTruncated
		session.mutex.Unlock()

		finished := "-"

		if finished_at.IsZero() == false {
			finished = finished_at.Format(time.RFC3339)
		}

		truncated := ""

		if log_truncated == true {
			truncated = ", LogTruncated: true"
		}

		return strings.Join([]string{
			fmt.Sprintf("shell.Stat: Session \"%s\": program \"%s\", arguments %s, status \"%s\"%s.", session_id, program, utils_shell.FormatArguments(arguments), status, utils_shell.FormatExitCode(status, exit_code)),
			fmt.Sprintf("shell.Stat: Workdir: \"%s\", TTY: %dx%d, Started: %s, Finished: %s, LastActivity: %s, Log: %d bytes, Dropped: %d bytes%s.", workdir, columns, rows, started_at.Format(time.RFC3339), finished, last_activity.Format(time.RFC3339), log_size, dropped, truncated),
		}, "\n"), nil

	} else {
		return "", fmt.Errorf("shell.Stat: Session \"%s\" not found.", session_id)
	}

}
