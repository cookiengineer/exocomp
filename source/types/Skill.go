package types

import "exocomp/encoding/yaml"
import "encoding/json"
import "errors"
import "strings"

type Skill struct {

	// frontmatter: agentskills.io Specification
	Name          string            `json:"name" yaml:"name"`
	Description   string            `json:"description" yaml:"description"`
	License       string            `json:"license" yaml:"license"`
	Compatibility string            `json:"compatibility" yaml:"compatibility"`
	Metadata      map[string]string `json:"metadata" yaml:"metadata"`
	AllowedTools  []string          `json:"allowed-tools" yaml:"allowed-tools"`

	// frontmatter: exocomp Specification
	AllowedPrograms []string        `json:"allowed-programs" yaml:"allowed-programs"`

	// body
	Body          string            `json:"body" yaml:"-"`

	// Internal Properties
	Scripts map[string]string `json:"-" yaml:"-"` // map[script]runtime

}

func ParseSkill(data []byte) (*Skill, error) {

	if len(data) > 2 && data[0] == '{' && data[len(data)-1] == '}' {

		skill := Skill{}
		err   := json.Unmarshal(data, &skill)

		if err == nil {
			return &skill, nil
		} else {
			return nil, err
		}

	} else {

		text := strings.TrimSpace(string(data))

		if strings.HasPrefix(text, "---\n") && strings.Contains(text, "\n---\n") {
			
			frontmatter := strings.TrimSpace(text[4:strings.Index(text, "\n---\n")])
			body        := strings.TrimSpace(text[strings.Index(text, "\n---\n")+5:])

			if len(frontmatter) > 0 && len(body) > 0 {

				skill := Skill{}
				err   := yaml.Unmarshal(data, &skill)

				if err == nil {

					skill.Body = body

					return &skill, nil

				} else {
					return nil, err
				}

			} else {
				return nil, errors.New("unsupported frontmatter syntax")
			}

		} else {
			return nil, errors.New("unsupported frontmatter syntax")
		}

	}

}
