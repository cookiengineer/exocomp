package vllm

import "strings"

type ModelConfigResponse struct {

	Model           string `json:"model"`
	MaxModelLength  int    `json:"max_model_len"`
	Tokenizer       string `json:"tokenizer"`
	TokenizerMode   string `json:"tokenizer_mode"`
	TrustRemoteCode bool   `json:"trust_remote_code"`
	DType           string `json:"dtype"`
	LoadFormat      string `json:"load_format"`
	ServedModelName string `json:"served_model_name"`

}

func (response *ModelConfigResponse) ContextLength(model_name string) int {

	result := int(0)

	if strings.Contains(response.Model, model_name) || strings.Contains(response.ServedModelName, model_name) {
		result = response.MaxModelLength
	}

	return result

}
