package agents

import "strings"

func render_prompt(name string, role string, template string) string {

	prompt := strings.TrimSpace(template)
	prompt = strings.ReplaceAll(prompt, "{{name}}", name)
	prompt = strings.ReplaceAll(prompt, "{{role}}", role)
	prompt = strings.TrimSpace(prompt)

	return prompt

}
