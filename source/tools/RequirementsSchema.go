package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Requirements.json
var requirementsJSON []byte

var RequirementsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(requirementsJSON, &schema)

	if err == nil {
		RequirementsSchema = schema
	} else {
		panic(err)
	}

}
