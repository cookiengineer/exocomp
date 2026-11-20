package programs

import "exocomp/schemas"
import _ "embed"
import "encoding/json"

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
