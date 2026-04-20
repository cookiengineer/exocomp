package agents

type AgentType string

const (
	AgentTypeArchitect  AgentType = "architect"
	AgentTypeCoder      AgentType = "coder"
	AgentTypeManager    AgentType = "manager"
	AgentTypeSummarizer AgentType = "summarizer"
	AgentTypeTester     AgentType = "tester"
	AgentTypeDefault    AgentType = "default"
)

func (agent_type AgentType) String() string {
	return string(agent_type)
}

