package schemas

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}
