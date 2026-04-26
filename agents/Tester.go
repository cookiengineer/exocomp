package agents

import "exocomp/types"
import _ "embed"
import "strings"

//go:embed Tester.txt
var tester_prompt []byte

func NewTester(config *types.Config) *types.Agent {

	name  := strings.TrimSpace(config.Name)
	model := strings.TrimSpace(config.Model)
	temp  := config.Temperature

	if name == "" {
		name = "Peanut Tester"
	}

	if model == "" {
		model = "qwen3-coder:30b"
	}

	if temp == 0.0 {
		temp = 0.3
	}

	return &types.Agent{
		Name:        name,
		Type:        "tester",
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

}
