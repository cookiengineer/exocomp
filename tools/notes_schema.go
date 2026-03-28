package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed notes_schema.json
var notes_json []byte

var notes_schema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(notes_json, &schema)

	if err == nil {
		notes_schema = schema
	} else {
		panic(err)
	}

}
