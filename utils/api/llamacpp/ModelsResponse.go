package llamacpp

import "slices"
import "strings"

type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

func (response *ModelsResponse) OwnedBy() string {

	owned_by := ""

	for _, model := range response.Data {

		if model.OwnedBy == "llamacpp" {
			owned_by = "llamacpp"
			break
		} else if model.OwnedBy == "library" {
			owned_by = "ollama"
			break
		} else if model.OwnedBy == "vllm" {
			owned_by = "vllm"
			break
		}

	}

	return owned_by

}


func (response *ModelsResponse) ContextLength(model_name string) int {

	result := int(0)

	for _, model := range response.Data {

		if strings.HasPrefix(model.ID, model_name) || slices.Contains(model.Aliases, model_name) {

			if model.OwnedBy == "llamacpp" {

				if model.Meta.ContextTrainingLength != 0 {
					result = model.Meta.ContextTrainingLength
					break
				}

			}

		}

	}

	return result

}
