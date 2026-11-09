package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Bugs.json
var bugsJSON []byte

var BugsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(bugsJSON, &schema)

	if err == nil {
		BugsSchema = schema
	} else {
		panic(err)
	}

}
