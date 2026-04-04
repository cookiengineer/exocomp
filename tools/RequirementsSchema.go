package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Requirements.json
var requirements_json []byte

var RequirementsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(requirements_json, &schema)

	if err == nil {
		RequirementsSchema = schema
	} else {
		panic(err)
	}

}
