package assets

import "exocomp/schemas"
import _ "embed"
import "encoding/json"

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
