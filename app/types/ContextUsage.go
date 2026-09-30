package types

type ContextUsage struct {
	Length           int     `json:"length" yaml:"length"`
	Tokens           int     `json:"tokens" yaml:"tokens"`
	PromptTokens     int     `json:"prompt_tokens" yaml:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens" yaml:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens" yaml:"total_tokens"`
	Cost             float64 `json:"cost" yaml:"cost"`
}
