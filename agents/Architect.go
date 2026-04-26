package agents

import "exocomp/schemas"
import "exocomp/types"
import _ "embed"
import "strings"

//go:embed Architect.txt
var architect_prompt []byte

func NewArchitect(config *types.Config) *types.Agent {

	name  := strings.TrimSpace(config.Name)
	model := strings.TrimSpace(config.Model)
	temp  := config.Temperature

	if name == "" {
		name = "Peanut Architect"
	}

	if model == "" {
		model = "gemma4:31b"
	}

	if temp == 0.0 {
		temp = 0.5
	}

	prompt   := renderPrompt(name, string(architect_prompt))
	messages := make([]*schemas.Message, 0)
	messages = append(messages, &schemas.Message{
		Role:    "system",
		Content: prompt,
	})

	return &types.Agent{
		Name:        name,
		Type:        "architect",
		Model:       model,
		Prompt:      prompt,
		Temperature: temp,
		Messages:    messages,
		Programs:    []string{},
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
		Sandbox:     config.Sandbox,
	}

}
