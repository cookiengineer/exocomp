package schemas

import "encoding/json"

type ToolCallFunction struct {
	Name         string          `json:"name" yaml:"name"`
	ArgumentsRaw json.RawMessage `json:"arguments" yaml:"arguments"` // some models return JSON string
	// Arguments map[string]interface{} `json:"arguments"`
}
