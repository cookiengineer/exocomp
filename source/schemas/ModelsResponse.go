package schemas

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

