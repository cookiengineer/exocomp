package tools

import "exocomp/agents"
// import "exocomp/types"
import "exocomp/utils"
import "fmt"
import "os"

type Agents struct {

	// TODO: running agents should be mapped to their states
	// TODO: Instruct agents to do a specific task
	// TODO: Start/Stop agents in a specific path, which must be in tool.Sandbox
	// TODO: agent == "manager" then allow to spawn other agents
	// TODO: agent == "coder" then allow to report back
	// TODO: agent == "tester" then allow to report back

	Sandbox   string
	workers   map[string]*agents.Agent
	processes map[string]*os.Process

}

func NewAgents(agent string, sandbox string, playground string) *Agents {

	agents := &Agents{
		Sandbox:   sandbox,
		workers:   make(map[string]*agents.Agent),
		processes: make(map[string]*os.Process),
	}

	return agents

}

func (tool *Agents) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "List" {

		return tool.List()

	} else if method == "Hire" {

		name,   ok1 := arguments["name"].(string)
		agent,  ok2 := arguments["agent"].(string)
		prompt, ok3 := arguments["prompt"].(string)

		if ok1 == true && ok2 == true && ok3 == true {
			return tool.Hire(utils.FormatAgentName(name), agent, utils.FormatMultiLine(prompt))
		} else {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Fire" {

		name, ok1 := arguments["name"].(string)

		if ok1 == true {
			return tool.Fire(utils.FormatAgentName(name))
		} else {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Quit" {

		message, ok1 := arguments["message"].(string)

		if ok1 == true {
			return tool.Quit(utils.FormatMultiLine(message))
		} else {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameters.")
		}

	} else {
		return "", fmt.Errorf("agents.%s: Invalid method.", method)
	}

}

func (tool *Agents) List() (string, error) {
	return "", fmt.Errorf("agents.List: Not implemented")
}

func (tool *Agents) Hire(name string, agent string, prompt string) (string, error) {
	return "", fmt.Errorf("agents.Hire: Not implemented")

}

func (tool *Agents) Fire(name string) (string, error) {
	return "", fmt.Errorf("agents.Fire: Not implemented")
}

func (tool *Agents) Quit(message string) (string, error) {
	return "", fmt.Errorf("agents.Quit: Not implemented")
}

