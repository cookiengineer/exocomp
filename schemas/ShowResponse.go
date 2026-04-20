package schemas

import "strings"

type ShowResponse struct {
	Capabilities []string       `json:"capabilities"`
	ModelInfo    map[string]any `json:"model_info"`
}

func (response *ShowResponse) ContextLength() int {

	result := int(0)

	if len(response.ModelInfo) > 0 {

		for key, raw_value := range response.ModelInfo {

			if strings.Contains(key, ".context_length") {

				switch tmp := raw_value.(type) {
				case int:
					result = tmp
				case int64:
					result = int(tmp)
				case float64:
					result = int(tmp)
				}

			}

		}

	}

	return result

}

func (response *ShowResponse) EmbeddingLength() int {

	result := int(0)

	if len(response.ModelInfo) > 0 {

		for key, raw_value := range response.ModelInfo {

			if strings.Contains(key, ".embedding_length") {

				switch tmp := raw_value.(type) {
				case int:
					result = tmp
				case int64:
					result = int(tmp)
				case float64:
					result = int(tmp)
				}

			}

		}

	}

	return result

}

