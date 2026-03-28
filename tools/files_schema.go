package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed files_schema.json
var files_json []byte

var files_schema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(files_json, &schema)

	if err == nil {
		files_schema = schema
	} else {
		panic(err)
	}

}
