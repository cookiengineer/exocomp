package types

import "exocomp/schemas"
import "exocomp/parsers/yaml"
import utils_fmt "exocomp/utils/fmt"
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
	Status          string             `json:"status,omitempty" yaml:"-"`
	StartedAt       schemas.Datetime   `json:"started-at,omitempty" yaml:"-"`
	FinishedAt      schemas.Datetime   `json:"finished-at,omitempty" yaml:"-"`
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

func (agent *Agent) HasTools() bool {

	if len(agent.AllowedTools) > 0 {
		return true
	}

	return false

}

func (agent *Agent) IsValid() bool {

	tmp_name        := utils_fmt.FormatAgentName(agent.Name)
	tmp_description := utils_fmt.FormatASCII(agent.Description)
	tmp_role        := utils_fmt.FormatAgentRole(agent.Role)
	tmp_model       := utils_fmt.FormatAgentModel(agent.Model)
	tmp_prompt      := utils_fmt.FormatASCII(agent.Prompt)
	tmp_sandbox     := utils_fmt.FormatFilePath(agent.Sandbox)

	if tmp_name == agent.Name &&
	   tmp_description == agent.Description &&
	   tmp_role == agent.Role &&
	   tmp_model == agent.Model &&
	   tmp_prompt == agent.Prompt &&
	   tmp_sandbox == agent.Sandbox {
		return true
	}

	return false

}
