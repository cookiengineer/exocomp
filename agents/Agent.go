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
	Type     AgentType
	Prompt   string
	Programs []string
	Tools    []string
}

func NewAgent(agent_type string) *Agent {

	switch agent_type {
	case "coder":

		return &Agent{
			Type:     AgentTypeCoder,
			Prompt:   string(coder_prompt),
			Programs: []string{"go", "gofmt", "gopls"},
			Tools:    []string{"bugs", "changelog", "features", "files", "programs"},
		}

	case "tester":

		return &Agent{
			Type:     AgentTypeTester,
			Prompt:   string(tester_prompt),
			Programs: []string{"go", "gofmt", "gopls"},
			Tools:    []string{"bugs", "changelog", "features", "files", "programs"},
		}

	case "manager":

		return &Agent{
			Type:     AgentTypeManager,
			Prompt:   string(manager_prompt),
			Programs: []string{},
			Tools:    []string{"agents", "features", "web"},
		}

	default:

		return &Agent{
			Type:     AgentType(""),
			Prompt:   string(default_prompt),
			Programs: []string{"go", "gofmt", "gopls"},
			Tools:    []string{"bugs", "changelog", "features", "files", "programs", "web"},
		}

	}

}

func (agent *Agent) GetPrompt() string {
	return string(agent.Prompt)
}
