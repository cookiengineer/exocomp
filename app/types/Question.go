package types

type Question struct {
	Type     string   `json:"type"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Multiple bool     `json:"multiple"`
	Answer   string   `json:"answer"`
}

func NewQuestion(arguments map[string]any) *Question {

	question := Question{
		Options: make([]string, 0),
	}

	if value, ok := arguments["question"].(string); ok == true {
		question.Question = value
	}

	if values, ok := arguments["options"].([]any); ok == true {

		for _, value := range values {

			if option, ok := value.(string); ok == true {
				question.Options = append(question.Options, option)
			}

		}

	} else if values, ok := arguments["options"].([]string); ok == true {
		question.Options = values
	}

	if value, ok := arguments["multiple"].(bool); ok == true {
		question.Multiple = value
	}

	if question.Question != "" && len(question.Options) > 0 {
		question.Type = "Choose"
	} else {
		question.Type = "Ask"
	}

	return &question

}

func (question *Question) IsChoice() bool {
	return question.Type == "Choose"
}
