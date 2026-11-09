package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Websites.json
var websitesJSON []byte

var WebsitesSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(websitesJSON, &schema)

	if err == nil {
		WebsitesSchema = schema
	} else {
		panic(err)
	}

}
