package agents

import "exocomp/schemas"
import "exocomp/types"
import _ "embed"
import "strings"

//go:embed Coder.txt
var coder_prompt []byte

func NewCoder(config *types.Config) *types.Agent {

	name  := strings.TrimSpace(config.Name)
	model := strings.TrimSpace(config.Model)
	temp  := config.Temperature

	if name == "" {
		name = "Peanut Coder"
	}

	if model == "" {
		model = "qwen3-coder:30b"
	}

	if temp == 0.0 {
		temp = 0.3
	}

	prompt   := renderPrompt(name, string(coder_prompt))
	messages := make([]*schemas.Message, 0)
	messages = append(messages, &schemas.Message{
		Role:    "system",
		Content: prompt,
	})

	return &types.Agent{
		Name:        name,
		Type:        "coder",
		Model:       model,
		Prompt:      prompt,
		Temperature: temp,
		Messages:    messages,
		Programs:    []string{"go", "gofmt", "gopls"},
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
		Sandbox:     config.Sandbox,
	}

}
