package ollama

type Message struct {
	Role    string `json:"role"`    // user || assistant || tool || system
	Content string `json:"content"`
}

