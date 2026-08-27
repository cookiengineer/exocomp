package types

type Requirement struct {
	Type          string `json:"type"`
	File          string `json:"file"`
	Symbol        string `json:"symbol"`
	Declaration   string `json:"declaration"`
	Behavior      string `json:"behavior"`
	IsImplemented bool   `json:"is_implemented"`
}
