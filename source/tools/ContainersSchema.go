package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Containers.json
var containersJSON []byte

var ContainersSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(containersJSON, &schema)

	if err == nil {
		ContainersSchema = schema
	} else {
		panic(err)
	}

}
