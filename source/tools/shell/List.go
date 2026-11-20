package shell

import utils_shell "exocomp/utils/shell"
import "fmt"
import "sort"
import "strings"
import "time"

func (tool *Shell) List() (string, error) {

	tool.mutex.Lock()
	sessions := make([]*shellSession, 0, len(tool.sessions))

	for _, session := range tool.sessions {
		sessions = append(sessions, session)
	}

	tool.mutex.Unlock()

	sort.Slice(sessions, func(a int, b int) bool {
		return sessions[a].StartedAt.Before(sessions[b].StartedAt)
	})

	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("shell.List: %d sessions.", len(sessions)))

	for _, session := range sessions {

		session.mutex.Lock()
		status := session.Status
		exit_code := session.ExitCode
		program := session.Program
		started_at := session.StartedAt
		last_activity := session.LastActivity
		log_size := session.logSize
		dropped := session.dropped
		session.mutex.Unlock()

		lines = append(lines, fmt.Sprintf("- Session: %s, Program: \"%s\", Status: %s%s, Started: %s, LastActivity: %s, Log: %d bytes, Dropped: %d bytes", session.ID, program, status, utils_shell.FormatExitCode(status, exit_code), started_at.Format(time.RFC3339), last_activity.Format(time.RFC3339), log_size, dropped))

	}

	return strings.Join(lines, "\n"), nil

}
