package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Agents.json
var agents_json []byte

var AgentsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(agents_json, &schema)

	if err == nil {
		AgentsSchema = schema
	} else {
		panic(err)
	}

}
