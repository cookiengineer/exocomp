package findings

type nucleiImportInfo struct {
	Name           string                     `json:"name"`
	Severity       string                     `json:"severity"`
	Description    string                     `json:"description"`
	Reference      []string                   `json:"reference"`
	Classification nucleiImportClassification `json:"classification"`
}
