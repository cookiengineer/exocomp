package types

import "exocomp/schemas"

type Agent struct {
	Name            string             `json:"name" yaml:"name"`
	Role            string             `json:"role" yaml:"role"`
	Model           string             `json:"model" yaml:"model"`
	Prompt          string             `json:"-" yaml:"prompt"` // Never expose Prompt in JSON API
	Temperature     float64            `json:"temperature" yaml:"temperature"`
	Messages        []*schemas.Message `json:"messages" yaml:"messages"`
	AllowedPrograms []string           `json:"allowed-programs" yaml:"allowed-programs"`
	AllowedTools    []string           `json:"allowed-tools" yaml:"allowed-tools"`
	Sandbox         string             `json:"sandbox" yaml:"sandbox"`
	ContextUsage    ContextUsage       `json:"context-usage" yaml:"context-usage"`
}

