package agents

type AgentType string

const (
	AgentTypeCoder   AgentType = "coder"
	AgentTypeTester  AgentType = "tester"
	AgentTypeManager AgentType = "manager"
)
