package agents

type AgentType string

const (
	AgentTypeCoder   AgentType = "coder"
	AgentTypeTester  AgentType = "tester"
	AgentTypeManager AgentType = "manager"
	AgentTypeNone    AgentType = ""
)

func (agent_type AgentType) String() string {
	return string(agent_type)
}
