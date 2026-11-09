package schemas

type ToolFunctionParameterProperty struct {
	Type        string                         `json:"type"`
	Description string                         `json:"description"`
	Enum        []string                       `json:"enum,omitempty"`
	Items       *ToolFunctionParameterProperty `json:"items,omitempty"`
}
