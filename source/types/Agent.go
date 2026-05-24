package types

import "exocomp/schemas"
import "exocomp/encoding/yaml"
import "encoding/json"

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
}

func ParseAgent(data []byte) (*Agent, error) {

	if len(data) > 2 && data[0] == '{' && data[len(data)-1] == '}' {

		agent := Agent{}
		err   := json.Unmarshal(data, &agent)

		if err == nil {
			return &agent, nil
		} else {
			return nil, err
		}

	} else {

		agent := Agent{}
		err   := yaml.Unmarshal(data, &agent)

		if err == nil {
			return &agent, nil
		} else {
			return nil, err
		}

	}

}
