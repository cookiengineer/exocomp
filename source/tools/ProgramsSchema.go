package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Programs.json
var programsJSON []byte

var ProgramsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(programsJSON, &schema)

	if err == nil {
		ProgramsSchema = schema
	} else {
		panic(err)
	}

}
