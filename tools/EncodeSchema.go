package tools

import "exocomp/schemas"
import _ "embed"
import "slices"

func EncodeSchema(allowed_tools []string) []*schemas.Tool {

	result := make([]*schemas.Tool, 0)

	for _, schema := range AgentsSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {
			result = append(result, &schema)
		}

	}

	for _, schema := range BugsSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {
			result = append(result, &schema)
		}

	}

	for _, schema := range ChangelogSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {
			result = append(result, &schema)
		}

	}

	for _, schema := range FilesSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {
			result = append(result, &schema)
		}

	}

	for _, schema := range ProgramsSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {
			result = append(result, &schema)
		}

	}

	for _, schema := range RequirementsSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {
			result = append(result, &schema)
		}

	}

	return result

}
