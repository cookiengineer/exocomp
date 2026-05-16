package ollama

type RunningModel struct {
	Name          string `json:"name"`
	Model         string `json:"model"`
	Size          int    `json:"size"`
	Digest        string `json:"digest"`
	// Details map[string]any `json:"details"`
	ExpiresAt     string `json:"expires_at"`
	SizeVRAM      int    `json:"size_vram"`
	ContextLength int    `json:"context_length"`
}
