package agents

import _ "embed"

//go:embed Agent.Coder.txt
var coder_prompt []byte

//go:embed Agent.Manager.txt
var manager_prompt []byte

//go:embed Agent.Tester.txt
var tester_prompt []byte

//go:embed Agent.prompt.txt
var default_prompt []byte

type Agent struct {
	Type   AgentType
	Prompt string
}

func NewAgent(agent_type string) *Agent {

	switch agent_type {
	case "coder":

		return &Agent{
			Type:   AgentTypeCoder,
			Prompt: string(coder_prompt),
		}

	case "tester":

		return &Agent{
			Type:   AgentTypeTester,
			Prompt: string(tester_prompt),
		}

	case "manager":

		return &Agent{
			Type:   AgentTypeManager,
			Prompt: string(manager_prompt),
		}

	default:

		return &Agent{
			Type:   AgentType(""),
			Prompt: string(default_prompt),
		}

	}

}

func (agent *Agent) GetPrompt() string {
	return string(agent.Prompt)
}
