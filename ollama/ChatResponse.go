package ollama

type ChatResponse struct {
	Message *Message `json:"message"`
}
