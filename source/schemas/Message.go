package schemas

type Message struct {
	Role      string     `json:"role" yaml:"role"`                                 // user || assistant || tool || system
	Content   string     `json:"content" yaml:"content"`                           // content
	ToolCalls []ToolCall `json:"tool_calls,omitempty" yaml:"tool-calls,omitempty"` // empty if content present
	ToolName  string     `json:"tool_name,omitempty" yaml:"tool-name,omitempty"`
}

