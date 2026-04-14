package agents

type AgentType string

const (
	AgentTypeArchitect  AgentType = "architect"
	AgentTypeCoder      AgentType = "coder"
	AgentTypeManager    AgentType = "manager"
	AgentTypeSummarizer AgentType = "summarizer"
	AgentTypeTester     AgentType = "tester"
	AgentTypeNone       AgentType = ""
)

func (agent_type AgentType) String() string {
	return string(agent_type)
}

