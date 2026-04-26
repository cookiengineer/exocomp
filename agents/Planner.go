package agents

import "exocomp/types"
import _ "embed"
import "strings"

//go:embed Planner.txt
var planner_prompt []byte

func NewPlanner(config *types.Config) *types.Agent {

	name  := strings.TrimSpace(config.Name)
	model := strings.TrimSpace(config.Model)
	temp  := config.Temperature

	if name == "" {
		name = "Peanut Hamper"
	}

	if model == "" {
		model = "gemma4:31b"
	}

	if temp == 0.0 {
		temp = 0.7
	}

	return &types.Agent{
		Name:        name,
		Type:        "planner",
		Model:       model,
		Prompt:      string(planner_prompt),
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

}
