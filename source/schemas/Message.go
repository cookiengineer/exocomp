package schemas

type Message struct {
	Role             string     `json:"role" yaml:"role"`
	Content          string     `json:"content" yaml:"content"`
	ReasoningContent string     `json:"reasoning_content,omitempty" yaml:"reasoning-content,omitempty"`
	ToolCallID       string     `json:"tool_call_id,omitempty" yaml:"tool-call-id,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty" yaml:"tool-calls,omitempty"`
	ToolName         string     `json:"tool_name,omitempty" yaml:"tool-name,omitempty"`
	Created          Datetime   `json:"created,omitempty" yaml:"created,omitempty"`
}
