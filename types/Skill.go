package types

type Skill struct {

	// Upstream agentskills.io specification
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	License       string            `json:"license"`
	Compatibility string            `json:"compatibility"`
	Metadata      map[string]string `json:"metadata"`
	AllowedTools  []string          `json:"allowed-tools"`

	// Internal Properties
	Scripts       map[string]string `json:"scripts"`

}
