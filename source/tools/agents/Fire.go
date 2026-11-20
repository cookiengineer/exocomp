package agents

import "fmt"

func (tool *Agents) Fire(name string) (string, error) {

	tool.mutex.Lock()
	process, running := tool.processes[name]
	agent, ok := tool.contents[name]

	if ok == true && running == true {
		agent.Status = "fired"
	}

	tool.mutex.Unlock()

	if running == true {

		process.cancel()
		<-process.done

		return fmt.Sprintf("agents.Fire: Agent \"%s\" fired.", name), nil

	} else {
		return "", fmt.Errorf("agents.Fire: Agent \"%s\" already quit!", name)
	}

}
