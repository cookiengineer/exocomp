package files

import "exocomp/schemas"
import _ "embed"
import "encoding/json"

//go:embed Files.json
var filesJSON []byte

var FilesSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(filesJSON, &schema)

	if err == nil {
		FilesSchema = schema
	} else {
		panic(err)
	}

}
