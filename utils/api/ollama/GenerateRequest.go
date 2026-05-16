package ollama

type GenerateRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	KeepAlive int    `json:"keep_alive"`
}
