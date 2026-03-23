package ollama

type ChatRequest struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
	Stream   bool       `json:"stream"`
}

