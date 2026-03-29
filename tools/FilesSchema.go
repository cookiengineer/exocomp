package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Files.json
var files_json []byte

var FilesSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(files_json, &schema)

	if err == nil {
		FilesSchema = schema
	} else {
		panic(err)
	}

}
