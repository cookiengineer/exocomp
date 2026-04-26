package agents

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

	return &types.Agent{
		Name:        name,
		Type:        "architect",
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

}
