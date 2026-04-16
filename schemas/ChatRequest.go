package schemas

type ChatRequest struct {
	Model       string     `json:"model"`
	Messages    []*Message `json:"messages"`
	Stream      bool       `json:"stream"`
	Temperature float64    `json:"temperature"`
	Tools       []*Tool    `json:"tools"`
	ToolChoice  string     `json:"tool_choice"`
}

