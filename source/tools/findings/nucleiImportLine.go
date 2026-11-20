package findings

type nucleiImportLine struct {
	TemplateID       string           `json:"template-id"`
	TemplatePath     string           `json:"template-path"`
	Info             nucleiImportInfo `json:"info"`
	Type             string           `json:"type"`
	Host             string           `json:"host"`
	URL              string           `json:"url"`
	MatchedAt        string           `json:"matched-at"`
	Request          string           `json:"request"`
	Response         string           `json:"response"`
	MatcherName      string           `json:"matcher-name"`
	ExtractedResults []string         `json:"extracted-results"`
}
