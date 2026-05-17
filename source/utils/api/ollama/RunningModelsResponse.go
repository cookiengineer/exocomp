package ollama

type RunningModelsResponse struct {
	Models []RunningModel `json:"models"`
}

func (response *RunningModelsResponse) ContextLength(model_name string) int {

	result := int(0)

	for _, model := range response.Models {

		if model.Name == model_name || model.Model == model_name {
			result = model.ContextLength
			break
		}

	}

	return result

}
