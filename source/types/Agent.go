package types

import "exocomp/schemas"

type Agent struct {
	Name            string             `json:"name"`
	Role            string             `json:"role"`
	Model           string             `json:"model"`
	Prompt          string             `json:"-"` // Never expose Prompt
	Temperature     float64            `json:"temperature"`
	Messages        []*schemas.Message `json:"messages"`
	AllowedPrograms []string           `json:"allowed-programs"`
	AllowedTools    []string           `json:"allowed-tools"`
	Sandbox         string             `json:"sandbox"`
	ContextUsage    ContextUsage       `json:"context-usage"`
}

