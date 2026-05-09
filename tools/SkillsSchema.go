package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Skills.json
var skills_json []byte

var SkillsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(skills_json, &schema)

	if err == nil {
		SkillsSchema = schema
	} else {
		panic(err)
	}

}
