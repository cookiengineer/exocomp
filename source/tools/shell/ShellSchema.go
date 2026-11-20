package shell

import "exocomp/schemas"
import _ "embed"
import "encoding/json"

//go:embed Shell.json
var shellJSON []byte

var ShellSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(shellJSON, &schema)

	if err == nil {
		ShellSchema = schema
	} else {
		panic(err)
	}

}
