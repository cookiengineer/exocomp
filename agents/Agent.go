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
	Type        AgentType
	Model       string
	Prompt      string
	Programs    []string
	Temperature float64
	Tools       []string
}

func NewAgent(agent_type string, agent_model string, agent_temperature float64) *Agent {

	// TODO: model
	// TODO: temperature

	switch agent_type {
	case "coder":

		return &Agent{
			Type:        AgentTypeCoder,
			Model:       "qwen3-coder:30b",
			Prompt:      string(coder_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: 0.3,
			Tools:       []string{"bugs", "changelog", "features", "files", "programs"},
		}

	case "tester":

		return &Agent{
			Type:        AgentTypeTester,
			Model:       "qwen3-coder:30b",
			Prompt:      string(tester_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: 0.3,
			Tools:       []string{"bugs", "changelog", "features", "files", "programs"},
		}

	case "manager":

		return &Agent{
			Type:        AgentTypeManager,
			Model:       "qwen3-coder:30b",
			Prompt:      string(manager_prompt),
			Programs:    []string{},
			Temperature: 0.7,
			Tools:       []string{"agents", "features", "web"},
		}

	default:

		return &Agent{
			Type:        AgentType(""),
			Model:       "qwen3-coder:30b",
			Prompt:      string(default_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: 0.5,
			Tools:       []string{"bugs", "changelog", "features", "files", "programs", "web"},
		}

	}

}

func (agent *Agent) GetPrompt() string {
	return string(agent.Prompt)
}
