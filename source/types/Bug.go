package types

type Bug struct {
	IsFixed     bool   `json:"is_fixed"`
	File        string `json:"file"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
}
