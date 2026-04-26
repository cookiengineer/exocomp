package agents

import "exocomp/schemas"
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

	prompt   := renderPrompt(name, string(summarizer_prompt))
	messages := make([]*schemas.Message, 0)
	messages = append(messages, &schemas.Message{
		Role:    "system",
		Content: prompt,
	})

	return &types.Agent{
		Name:        name,
		Type:        "summarizer",
		Model:       model,
		Prompt:      prompt,
		Temperature: temp,
		Messages:    messages,
		Programs:    []string{"go", "gopls"},
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
		Sandbox:     config.Sandbox,
	}

}
