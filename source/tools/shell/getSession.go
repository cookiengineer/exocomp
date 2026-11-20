package shell

func (tool *Shell) getSession(session_id string) *shellSession {

	tool.mutex.Lock()
	session, ok := tool.sessions[session_id]
	tool.mutex.Unlock()

	if ok == true {
		return session
	} else {
		return nil
	}

}
