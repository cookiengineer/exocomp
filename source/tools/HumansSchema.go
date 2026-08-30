package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Humans.json
var humans_json []byte

var HumansSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(humans_json, &schema)

	if err == nil {
		HumansSchema = schema
	} else {
		panic(err)
	}

}
