package tools

import "exocomp/schemas"
import _ "embed"
import "slices"

func EncodeSchema(allowed_tools []string) []schemas.Tool {

	result := make([]schemas.Tool, 0)

	if slices.Contains(allowed_tools, "agents") {

		// TODO

	} else if slices.Contains(allowed_tools, "bugs") {

		// TODO

	} else if slices.Contains(allowed_tools, "changelog") {

		for _, schema := range ChangelogSchema {
			result = append(result, schema)
		}

	} else if slices.Contains(allowed_tools, "features") {

		// TODO

	} else if slices.Contains(allowed_tools, "files") {

		for _, schema := range FilesSchema {
			result = append(result, schema)
		}

	} else if slices.Contains(allowed_tools, "programs") {

		for _, schema := range ProgramsSchema {
			result = append(result, schema)
		}

	} else if slices.Contains(allowed_tools, "web") {

		// TODO

	}

	return result

}
