package agents

import "exocomp/schemas"
import "exocomp/types"
import _ "embed"
import "strings"

//go:embed WebScanner.txt
var webscanner_prompt []byte

func NewWebScanner(config *types.Config) *types.Agent {

	name  := strings.TrimSpace(config.Name)
	model := strings.TrimSpace(config.Model)
	temp  := config.Temperature

	if name == "" {
		name = "Peanut Webscanner"
	}

	if model == "" {
		model = "qwen3-coder:30b"
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
		Name:            name,
		Role:            "webscanner",
		Model:           model,
		Prompt:          prompt,
		Temperature:     temp,
		Messages:        messages,
		AllowedPrograms: []string{
			"curl",
			"amass",
			"asnmap",
			"go",
			"httpx",
			"naabu",
			"nuclei",
			"subfinder",
		},
		AllowedTools:    []string{
			// No agents hiring
			"agents.Quit",
			// No bugs
			// No changelog
			"files.List",
			"files.Read",
			"files.Stat",
			"files.Write",
			"programs.List",
			"programs.Execute",
			// No requirements
		},
		Sandbox: config.Sandbox,
	}

}
