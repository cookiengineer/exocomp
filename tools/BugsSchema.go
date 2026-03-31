package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Bugs.json
var bugs_json []byte

var BugsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(bugs_json, &schema)

	if err == nil {
		BugsSchema = schema
	} else {
		panic(err)
	}

}
