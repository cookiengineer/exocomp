package humans

import "exocomp/schemas"
import _ "embed"
import "encoding/json"

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
