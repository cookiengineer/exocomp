package types

import "strings"

type Agent struct {
	Name        string
	Type        string
	Model       string
	Prompt      string
	Programs    []string
	Temperature float64
	Tools       []string
}

func (agent *Agent) GetPrompt() string {

	prompt := agent.Prompt
	prompt = strings.ReplaceAll(prompt, "{{NAME}}", agent.Name)
	prompt = strings.TrimSpace(prompt)

	return prompt

}
