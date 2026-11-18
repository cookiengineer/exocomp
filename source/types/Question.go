package types

type Question struct {
	Type         string       `json:"type"`
	Agent        string       `json:"agent,omitempty"`
	Question     string       `json:"question"`
	Options      []string     `json:"options"`
	Multiple     bool         `json:"multiple"`
	Capabilities []Capability `json:"capabilities,omitempty"`
	Answer       string       `json:"answer"`
}
