package schemas

import "encoding/json"

type ToolFunctionParameterProperty struct {
	Type        string                                   `json:"type"`
	Description string                                   `json:"description"`
	Enum        []string                                 `json:"enum,omitempty"`
	Items       *ToolFunctionParameterProperty           `json:"items,omitempty"`
	Properties  map[string]ToolFunctionParameterProperty `json:"properties,omitempty"`
	Required    []string                                 `json:"required,omitempty"`
	Default     json.RawMessage                          `json:"default,omitempty"`
	Minimum     *float64                                 `json:"minimum,omitempty"`
	Maximum     *float64                                 `json:"maximum,omitempty"`
}
