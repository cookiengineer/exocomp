package tools

type Agents struct {
	// TODO: running agents should be mapped to their states
	// TODO: Instruct agents to do a specific task
	// TODO: Start/Stop agents in a specific path, which must be in tool.Sandbox
	Sandbox string
	running map[string]AgentState
}

func NewAgents(agent string, sandbox string) *Agents {

	// TODO: agent == "manager" then allow to spawn other agents
	// TODO: agent == "coder" then allow to report back
	// TODO: agent == "tester" then allow to report back

}
