package agents

import _ "embed"
import "strings"

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

func NewAgent(agent_name string, agent_type string, agent_model string, agent_temperature float64) *Agent {

	if agent_type == "architect" {

		name := strings.TrimSpace(agent_name)

		if name == "" {
			name = "Peanut Architect"
		}


		return &Agent{
			Name:        name,
			Type:        AgentTypeArchitect,
			Model:       "qwen3-coder:30b",
			Prompt:      string(architect_prompt),
			Programs:    []string{},
			Temperature: 0.5,
			Tools:       []string{
				"agents.List",
				"agents.Hire",
				"agents.Fire",
				"agents.Quit",
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

	} else if agent_type == "coder" {

		name := strings.TrimSpace(agent_name)

		if name == "" {
			name = "Peanut Coder"
		}

		return &Agent{
			Name:        name,
			Type:        AgentTypeCoder,
			Model:       "qwen3-coder:30b",
			Prompt:      string(coder_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: 0.3,
			Tools:       []string{
				"agents.Quit",
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

	} else if agent_type == "tester" {

		name := strings.TrimSpace(agent_name)

		if name == "" {
			name = "Peanut Tester"
		}

		return &Agent{
			Name:        name,
			Type:        AgentTypeTester,
			Model:       "qwen3-coder:30b",
			Prompt:      string(tester_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: 0.3,
			Tools:       []string{
				"agents.Quit",
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

	} else if agent_type == "manager" {

		name := strings.TrimSpace(agent_name)

		if name == "" {
			name = "Peanut Manager"
		}

		return &Agent{
			Name:        name,
			Type:        AgentTypeManager,
			Model:       "qwen3-coder:30b",
			Prompt:      string(manager_prompt),
			Programs:    []string{},
			Temperature: 0.7,
			Tools:       []string{
				"agents.List",
				"agents.Hire",
				"agents.Fire",
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

	} else {

		name := strings.TrimSpace(agent_name)

		if name == "" {
			name = "Peanut Hamper"
		}

		return &Agent{
			Name:        name,
			Type:        AgentTypeNone,
			Model:       "qwen3-coder:30b",
			Prompt:      string(none_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: 0.5,
			Tools:       []string{
				"agents.List",
				"agents.Hire",
				"agents.Fire",
				"agents.Quit",
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
