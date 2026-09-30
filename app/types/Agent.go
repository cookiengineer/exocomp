package types

import "exocomp-app/schemas"

type Agent struct {
	Name            string             `json:"name" yaml:"name"`
	Description     string             `json:"description" yaml:"description"`
	Role            string             `json:"role" yaml:"role"`
	Model           string             `json:"model" yaml:"model"`
	Prompt          string             `json:"prompt" yaml:"prompt"`
	Temperature     float64            `json:"temperature" yaml:"temperature"`
	Messages        []*schemas.Message `json:"messages" yaml:"messages"`
	AllowedPrograms []string           `json:"allowed_programs" yaml:"allowed-programs"`
	AllowedTools    []string           `json:"allowed_tools" yaml:"allowed-tools"`
	Sandbox         string             `json:"sandbox" yaml:"-"`
	ContextUsage    ContextUsage       `json:"context-usage" yaml:"-"`
	Status          string             `json:"status,omitempty" yaml:"-"`
	StartedAt       schemas.Datetime   `json:"started-at,omitempty" yaml:"-"`
	FinishedAt      schemas.Datetime   `json:"finished-at,omitempty" yaml:"-"`
}

func NewAgent() *Agent {

	agent := Agent{
		Messages:        make([]*schemas.Message, 0),
		AllowedPrograms: make([]string, 0),
		AllowedTools:    make([]string, 0),
	}

	return &agent

}
