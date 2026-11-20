package skills

import "exocomp/schemas"
import _ "embed"
import "encoding/json"

//go:embed Skills.json
var skillsJSON []byte

var SkillsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(skillsJSON, &schema)

	if err == nil {
		SkillsSchema = schema
	} else {
		panic(err)
	}

}
