package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Websites.json
var websites_json []byte

var WebsitesSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(websites_json, &schema)

	if err == nil {
		WebsitesSchema = schema
	} else {
		panic(err)
	}

}
