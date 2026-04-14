package schemas

type ModelsResponse struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}
