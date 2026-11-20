package shell

import "time"

func (tool *Shell) cleanup() {

	now := time.Now()
	expired := make([]*shellSession, 0)

	tool.mutex.Lock()

	for session_id, session := range tool.sessions {

		session.mutex.Lock()
		status := session.Status
		last_activity := session.LastActivity
		started_at := session.StartedAt
		finished_at := session.FinishedAt
		session.mutex.Unlock()

		if (status == "starting" || status == "running") && (now.Sub(last_activity) > tool.IdleTimeout || now.Sub(started_at) > tool.Timeout) {
			delete(tool.sessions, session_id)
			expired = append(expired, session)
		} else if (status == "exited" || status == "failed" || status == "timed_out") && now.Sub(finished_at) > shellExitedRetention {
			delete(tool.sessions, session_id)
			expired = append(expired, session)
		}

	}

	tool.mutex.Unlock()

	for _, session := range expired {
		session.timeout()
		session.close(true)
	}

}
