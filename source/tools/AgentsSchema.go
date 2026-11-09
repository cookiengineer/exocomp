package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Agents.json
var agentsJSON []byte

var AgentsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(agentsJSON, &schema)

	if err == nil {
		AgentsSchema = schema
	} else {
		panic(err)
	}

}
