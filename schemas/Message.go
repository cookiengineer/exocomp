package schemas

type Message struct {
	Role      string     `json:"role"`                 // user || assistant || tool || system
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"` // empty if content present
	ToolName  string     `json:"tool_name,omitempty"`
}

