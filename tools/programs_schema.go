package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed programs_schema.json
var programs_json []byte

var programs_schema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(programs_json, &schema)

	if err == nil {
		programs_schema = schema
	} else {
		panic(err)
	}

}
