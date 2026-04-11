package agents

import _ "embed"
import "strings"

//go:embed Agent.Architect.txt
var architect_prompt []byte

//go:embed Agent.Coder.txt
var coder_prompt []byte

//go:embed Agent.Manager.txt
var manager_prompt []byte

//go:embed Agent.Summarizer.txt
var summarizer_prompt []byte

//go:embed Agent.Tester.txt
var tester_prompt []byte

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

		name  := strings.TrimSpace(agent_name)
		model := strings.TrimSpace(agent_model)
		temp  := agent_temperature

		if name == "" {
			name = "Peanut Architect"
		}

		if model == "" {
			model = "gemma4:31b"
		}

		if temp == 0.0 {
			temp = 0.5
		}

		return &Agent{
			Name:        name,
			Type:        AgentTypeArchitect,
			Model:       model,
			Prompt:      string(architect_prompt),
			Programs:    []string{},
			Temperature: temp,
			Tools:       []string{
				// No agents.List
				// No agents.Hire
				// No agents.Fire
				"agents.Quit",
				// No bugs.Add
				// No bugs.Fix
				"bugs.List",
				"bugs.Search",
				// No changelog
				"files.List",
				"files.Read",
				"files.Stat",
				// No files.Write
				// No programs
				"requirements.DefineFunc",
				"requirements.DefineStruct",
				"requirements.DefineTest",
				"requirements.List",
				"requirements.Search",
			},
		}

	} else if agent_type == "coder" {

		name  := strings.TrimSpace(agent_name)
		model := strings.TrimSpace(agent_model)
		temp  := agent_temperature

		if name == "" {
			name = "Peanut Coder"
		}

		if model == "" {
			model = "qwen3-coder:30b"
		}

		if temp == 0.0 {
			temp = 0.3
		}

		return &Agent{
			Name:        name,
			Type:        AgentTypeCoder,
			Model:       model,
			Prompt:      string(coder_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: temp,
			Tools:       []string{
				// No agents.List
				// No agents.Hire
				// No agents.Fire
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
				// No requirements.DefineFunc
				// No requirements.DefineStruct
				// No requirements.DefineTest
				"requirements.List",
				"requirements.Search",
			},
		}

	} else if agent_type == "researcher" {

		// TODO: Implement researcher

		return nil

	} else if agent_type == "summarizer" {

		name  := strings.TrimSpace(agent_name)
		model := strings.TrimSpace(agent_model)
		temp  := agent_temperature

		if name == "" {
			name = "Peanut Summarizer"
		}

		if model == "" {
			model = "gemma4:31b"
		}

		if temp == 0.0 {
			temp = 0.3
		}

		return &Agent{
			Name:        name,
			Type:        AgentTypeSummarizer,
			Model:       model,
			Prompt:      string(tester_prompt),
			Programs:    []string{"go", "gopls"},
			Temperature: temp,
			Tools:       []string{
				// No agents.List
				// No agents.Hire
				// No agents.Fire
				"agents.Quit",
				// No bugs.Add
				// No bugs.Fix
				"bugs.List",
				"bugs.Search",
				// No changelog
				"files.List",
				"files.Read",
				"files.Stat",
				// No files.Write
				"programs.List",
				"programs.Execute",
				// No requirements.DefineFunc
				// No requirements.DefineStruct
				// No requirements.DefineTest
				"requirements.List",
				"requirements.Search",
			},
		}

	} else if agent_type == "tester" {

		name  := strings.TrimSpace(agent_name)
		model := strings.TrimSpace(agent_model)
		temp  := agent_temperature

		if name == "" {
			name = "Peanut Tester"
		}

		if model == "" {
			model = "qwen3-coder:30b"
		}

		if temp == 0.0 {
			temp = 0.3
		}

		return &Agent{
			Name:        name,
			Type:        AgentTypeTester,
			Model:       model,
			Prompt:      string(tester_prompt),
			Programs:    []string{"go", "gofmt", "gopls"},
			Temperature: temp,
			Tools:       []string{
				// No agents.List
				// No agents.Hire
				// No agents.Fire
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
				// No requirements.DefineFunc
				// No requirements.DefineStruct
				// No requirements.DefineTest
				"requirements.List",
				"requirements.Search",
			},
		}

	} else if agent_type == "manager" {

		name  := strings.TrimSpace(agent_name)
		model := strings.TrimSpace(agent_model)
		temp  := agent_temperature

		if name == "" {
			name = "Peanut Manager"
		}

		if model == "" {
			model = "gemma4:31b"
		}

		if temp == 0.0 {
			temp = 0.7
		}

		return &Agent{
			Name:        name,
			Type:        AgentTypeManager,
			Model:       model,
			Prompt:      string(manager_prompt),
			Programs:    []string{},
			Temperature: temp,
			Tools:       []string{
				"agents.List",
				"agents.Hire",
				"agents.Fire",
				// No agents.Quit
				// No bugs.Add
				// No bugs.Fix
				// No bugs.List
				// No bugs.Search
				"files.List",
				"files.Read",
				"files.Stat",
				// No files.Write
				// No programs
				// No requirements.DefineFunc
				// No requirements.DefineStruct
				// No requirements.DefineTest
				// No requirements.List
				// No requirements.Search
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
				"requirements.DefineFunc",
				"requirements.DefineStruct",
				"requirements.DefineTest",
				"requirements.List",
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
