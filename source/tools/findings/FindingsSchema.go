package findings

import "exocomp/schemas"
import _ "embed"
import "encoding/json"

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
