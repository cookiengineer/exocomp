package schemas

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	Parameters  map[string]ToolFunctionParameter `json:"parameters"`
}

type ToolFunctionParameter struct {
	Type       string                                   `json:"type"`
	Properties map[string]ToolFunctionParameterProperty `json:"properties"`
	Required   []string                                 `json:"required"`
}

type ToolFunctionParameterProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

