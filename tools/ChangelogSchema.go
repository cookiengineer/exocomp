package tools

import "exocomp/schemas"
import "encoding/json"
import _ "embed"

//go:embed Changelog.json
var changelog_json []byte

var ChangelogSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err    := json.Unmarshal(changelog_json, &schema)

	if err == nil {
		ChangelogSchema = schema
	} else {
		panic(err)
	}

}
