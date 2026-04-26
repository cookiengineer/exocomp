package agents

import "exocomp/types"
import _ "embed"
import "strings"

//go:embed Summarizer.txt
var summarizer_prompt []byte

func NewSummarizer(config *types.Config) *types.Agent {

	name  := strings.TrimSpace(config.Name)
	model := strings.TrimSpace(config.Model)
	temp  := config.Temperature

	if name == "" {
		name = "Peanut Summarizer"
	}

	if model == "" {
		model = "gemma4:31b"
	}

	if temp == 0.0 {
		temp = 0.3
	}

	return &types.Agent{
		Name:        name,
		Type:        "summarizer",
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

}
