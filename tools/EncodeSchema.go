package tools

import "exocomp/schemas"
import _ "embed"
import "slices"

func EncodeSchema(allowed_tools []string) []schemas.Tool {

	result := make([]schemas.Tool, 0)

	if slices.Contains(allowed_tools, "bugs") {
		// TODO
	} else if slices.Contains(allowed_tools, "files") {

		for _, schema := range files_schema {
			result = append(result, schema)
		}

	} else if slices.Contains(allowed_tools, "notes") {

		for _, schema := range notes_schema {
			result = append(result, schema)
		}

	} else if slices.Contains(allowed_tools, "programs") {

		for _, schema := range programs_schema {
			result = append(result, schema)
		}

	} else if slices.Contains(allowed_tools, "features") {
		// TODO
	}

	return result

}
