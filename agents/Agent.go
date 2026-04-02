package agents

import _ "embed"

//go:embed Agent.Architect.txt
var architect_prompt []byte

//go:embed Agent.Coder.txt
var coder_prompt []byte

//go:embed Agent.Tester.txt
var tester_prompt []byte

//go:embed Agent.Manager.txt
var manager_prompt []byte

//go:embed Agent.txt
var none_prompt []byte

type Agent struct {
	Name        string
	Type        AgentType
	Model       string
	Prompt      string
	Programs    []string
	Temperature float64
	Tools       []string
}

func NewAgent(agent_type string, agent_model string, agent_temperature float64) *Agent {

	switch agent_type {
	case "architect":

		return &Agent{
			Name:        "Peanut Architect",
			Type:        AgentTypeArchitect,
			Model:       "qwen3-coder:30b",
			Prompt:      string(architect_prompt),
			Programs:    []string{},
			Temperature: 0.5,
			Tools:       []string{
				"agents.List",
				"agents.Message",
				"agents.Report",
				// No agents.Start
				// No agents.Stop
				// No bugs.Add
				// No bugs.Fix
				"bugs.List",
				"bugs.Search",
				"files.List",
				"files.Read",
				"files.Stat",
				// No files.Write
				"requirements.Add",
				"requirements.List",
				"requirements.Remove",
				"requirements.Search",
			},
		}

	case "coder":

		return &Agent{
			Name:        "Peanut Coder",
			Type:        AgentTypeCoder,
			Model:       "qwen3-coder:30b",
			Prompt:      string(coder_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: 0.3,
			Tools:       []string{
				"agents.List",
				"agents.Message",
				"agents.Report",
				// No agents.Start
				// No agents.Stop
				// No bugs.Add
				"bugs.Fix",
				"bugs.List",
				"bugs.Search",
				"changelog.Add",
				"changelog.Change",
				"changelog.Deprecate",
				"changelog.Fix",
				"changelog.Remove",
				"changelog.Search",
				"files.List",
				"files.Read",
				"files.Stat",
				"files.Write",
				"programs.List",
				"programs.Execute",
				// No requirements.Add
				"requirements.List",
				// No requirements.Remove
				"requirements.Search",
			},
		}

	case "tester":

		return &Agent{
			Name:        "Peanut Tester",
			Type:        AgentTypeTester,
			Model:       "qwen3-coder:30b",
			Prompt:      string(tester_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: 0.3,
			Tools:       []string{
				"agents.List",
				"agents.Message",
				"agents.Report",
				// No agents.Start
				// No agents.Stop
				"bugs.Add",
				// No bugs.Fix
				"bugs.List",
				"bugs.Search",
				// No changelog
				"files.List",
				"files.Read",
				"files.Stat",
				"files.Write",
				"programs.List",
				"programs.Execute",
				// No requirements.Add
				"requirements.List",
				// No requirements.Remove
				"requirements.Search",
			},
		}

	case "manager":

		return &Agent{
			Name:        "Peanut Manager",
			Type:        AgentTypeManager,
			Model:       "qwen3-coder:30b",
			Prompt:      string(manager_prompt),
			Programs:    []string{},
			Temperature: 0.7,
			Tools:       []string{
				"agents.List",
				"agents.ListMessages",
				"agents.ListReports",
				"agents.Message",
				"agents.Report",
				"agents.Start",
				"agents.Stop",
				// No bugs.Add
				// No bugs.Fix
				"bugs.List",
				"bugs.Search",
				"files.List",
				"files.Read",
				"files.Stat",
				// No files.Write
				"requirements.Add",
				"requirements.List",
				"requirements.Remove",
				"requirements.Search",
			},
		}

	default:

		return &Agent{
			Name:        "Peanut Hamper",
			Type:        AgentTypeNone,
			Model:       "qwen3-coder:30b",
			Prompt:      string(none_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: 0.5,
			Tools:       []string{
				// No agents
				"bugs.Add",
				"bugs.Fix",
				"bugs.List",
				"bugs.Search",
				"changelog.Add",
				"changelog.Change",
				"changelog.Deprecate",
				"changelog.Fix",
				"changelog.Remove",
				"changelog.Search",
				"files.List",
				"files.Read",
				"files.Stat",
				"files.Write",
				"programs.List",
				"programs.Execute",
				"requirements.Add",
				"requirements.List",
				"requirements.Remove",
				"requirements.Search",
			},
		}

	}

}

func (agent *Agent) GetPrompt() string {

	prompt := agent.Prompt

	prompt = strings.ReplaceAll(prompt, "{{NAME}}", agent.Name)

	return strings.TrimSpace(prompt)

}
