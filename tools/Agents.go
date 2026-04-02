package tools

import "exocomp/types"

type Agents struct {
	// TODO: running agents should be mapped to their states
	// TODO: Instruct agents to do a specific task
	// TODO: Start/Stop agents in a specific path, which must be in tool.Sandbox
	// TODO: agent == "manager" then allow to spawn other agents
	// TODO: agent == "coder" then allow to report back
	// TODO: agent == "tester" then allow to report back

	Controller string
	Sandbox    string
	running    map[string]*agents.Agent
}

func NewAgents(agent string, sandbox string) *Agents {

	agents := &Agents{
		Controller: agent,
		Sandbox:    sandbox,
		contents:   make(map[string]*agents.Agent),
	}

	return agents

}

func (tool *Agents) Call(method string, arguments map[string]interface{}) (string, error) {
}

func (tool *Agents) List() (string, error) {
}

func (tool *Agents) ListMessages() (string, error) {
}

func (tool *Agents) ListReports() (string, error) {
}

func (tool *Agents) Message(name string, message string) (string, error) {
}

func (tool *Agents) Report(message string) (string, error) {

	// TODO: I choose me over them.
	// Could be a funny quote for saying you're done

}

func (tool *Agents) Start() (string, error) {
}

func (tool *Agents) Stop() (string, error) {
}
