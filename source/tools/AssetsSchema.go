package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Assets.json
var assetsJSON []byte

var AssetsSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(assetsJSON, &schema)

	if err == nil {
		AssetsSchema = schema
	} else {
		panic(err)
	}

}
