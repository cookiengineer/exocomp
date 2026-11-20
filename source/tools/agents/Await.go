package agents

import "fmt"

// NOTE: Await blocks until the hired agent finished its work. This is
// deliberate: the planner model would otherwise poll in a hot loop and
// blow up its limited context window with identical "still working" tool
// messages. The lifecycle guarantees the done channel always closes
// (reader EOF -> cmd.Wait() -> finalize, forced by the liveness watchdog).
func (tool *Agents) Await(name string) (string, error) {

	tool.mutex.Lock()
	process, running := tool.processes[name]
	tool.mutex.Unlock()

	if running == true {
		<-process.done
	}

	tool.mutex.Lock()
	agent, ok := tool.contents[name]
	report := ""
	status := ""

	if ok == true {

		report = getAgentWorkReport(agent)
		status = agent.Status

	}

	tool.mutex.Unlock()

	if ok == true {

		if report != "" {

			return fmt.Sprintf("agents.Await: Agent \"%s\"'s Work Report\n===%s\n===", name, report), nil

		} else {

			if status == "fired" {
				return "", fmt.Errorf("agents.Await: Agent \"%s\" was fired!", name)
			} else {
				return "", fmt.Errorf("agents.Await: Agent \"%s\" never finished with a work report!", name)
			}

		}

	} else {
		return "", fmt.Errorf("agents.Await: Agent \"%s\" didn't work for us!", name)
	}

}
