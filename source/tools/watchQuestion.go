package tools

func watchQuestion(tool *Humans, question string, state *questionState) {

	<-state.ctx.Done()

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	_, still_waiting := tool.states[question]

	if still_waiting == true {

		delete(tool.states, question)
		close(state.done)

	}

}
