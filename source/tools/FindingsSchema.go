package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Findings.json
var findingsJSON []byte

var FindingsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(findingsJSON, &schema)

	if err == nil {
		FindingsSchema = schema
	} else {
		panic(err)
	}

}
