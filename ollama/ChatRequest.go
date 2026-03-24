package ollama

type ChatRequest struct {
	Model       string     `json:"model"`
	Messages    []*Message `json:"messages"`
	Stream      bool       `json:"stream"`
	Temperature float32    `json:"temperature"`
}

