package containers

import "exocomp/schemas"
import _ "embed"
import "encoding/json"

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
