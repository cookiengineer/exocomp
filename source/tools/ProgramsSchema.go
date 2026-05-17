package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Programs.json
var programs_json []byte

var ProgramsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(programs_json, &schema)

	if err == nil {
		ProgramsSchema = schema
	} else {
		panic(err)
	}

}
