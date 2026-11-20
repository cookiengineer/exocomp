package changelog

import "exocomp/schemas"
import _ "embed"
import "encoding/json"

//go:embed Changelog.json
var changelogJSON []byte

var ChangelogSchema []schemas.Tool

func init() {

	schema := make([]schemas.Tool, 0)
	err := json.Unmarshal(changelogJSON, &schema)

	if err == nil {
		ChangelogSchema = schema
	} else {
		panic(err)
	}

}
