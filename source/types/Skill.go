package types

type Skill struct {

	// agentskills.io Specification
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	License       string            `json:"license"`
	Compatibility string            `json:"compatibility"`
	Metadata      map[string]string `json:"metadata"`
	AllowedTools  []string          `json:"allowed-tools"`
	Body          string            `json:"body"`

	// exocomp Specification
	AllowedPrograms []string `json:"allowed-programs"`

	// Internal Properties
	Scripts map[string]string `json:"scripts"` // map[script]runtime

}
