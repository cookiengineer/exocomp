package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Questions.json
var questions_json []byte

var QuestionsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(questions_json, &schema)

	if err == nil {
		QuestionsSchema = schema
	} else {
		panic(err)
	}

}
