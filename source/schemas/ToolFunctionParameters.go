package schemas

type ToolFunctionParameters struct {
	Type       string                                   `json:"type"`
	Properties map[string]ToolFunctionParameterProperty `json:"properties"`
	Required   []string                                 `json:"required"`
}
