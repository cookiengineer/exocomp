package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Humans.json
var humansJSON []byte

var HumansSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(humansJSON, &schema)

	if err == nil {
		HumansSchema = schema
	} else {
		panic(err)
	}

}
