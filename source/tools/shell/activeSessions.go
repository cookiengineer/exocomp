package shell

func (tool *Shell) activeSessions() int {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	count := 0

	for _, session := range tool.sessions {

		session.mutex.Lock()

		if session.Status == "starting" || session.Status == "running" {
			count++
		}

		session.mutex.Unlock()

	}

	return count

}
