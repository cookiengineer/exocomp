package agents

import "strings"

func renderPrompt(name string, template string) string {

	prompt := strings.TrimSpace(template)
	prompt = strings.ReplaceAll(prompt, "{{NAME}}", name)
	prompt = strings.TrimSpace(prompt)

	return prompt

}
