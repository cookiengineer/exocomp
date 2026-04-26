package tools

import "exocomp/schemas"
import "exocomp/types"
import _ "embed"
import net_url "net/url"
import "slices"

func Toolset(playground string, sandbox string, url *net_url.URL, allowed_programs []string, allowed_tools []string) (map[string][]schemas.Tool, map[string]types.Tool) {

	result_schemas := make(map[string][]schemas.Tool, 0)
	result_tools   := make(map[string]types.Tool, 0)

	for _, schema := range AgentsSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {

			_, ok1 := result_schemas["agents"]

			if ok1 == false {
				result_schemas["agents"] = make([]schemas.Tool, 0)
			}

			result_schemas["agents"] = append(result_schemas["agents"], schema)

			_, ok2 := result_tools["agents"]

			if ok2 == false {
				result_tools["agents"] = NewAgents(playground, sandbox, url)
			}

		}

	}

	for _, schema := range BugsSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {

			_, ok1 := result_schemas["bugs"]

			if ok1 == false {
				result_schemas["bugs"] = make([]schemas.Tool, 0)
			}

			result_schemas["bugs"] = append(result_schemas["bugs"], schema)

			_, ok2 := result_tools["bugs"]

			if ok2 == false {
				result_tools["bugs"] = NewBugs(playground, sandbox)
			}

		}

	}

	for _, schema := range ChangelogSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {

			_, ok1 := result_schemas["changelog"]

			if ok1 == false {
				result_schemas["changelog"] = make([]schemas.Tool, 0)
			}

			result_schemas["changelog"] = append(result_schemas["changelog"], schema)

			_, ok2 := result_tools["changelog"]

			if ok2 == false {
				result_tools["changelog"] = NewChangelog(playground, sandbox)
			}

		}

	}

	for _, schema := range FilesSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {

			_, ok1 := result_schemas["files"]

			if ok1 == false {
				result_schemas["files"] = make([]schemas.Tool, 0)
			}

			result_schemas["files"] = append(result_schemas["files"], schema)

			_, ok2 := result_tools["files"]

			if ok2 == false {
				result_tools["files"] = NewFiles(playground, sandbox)
			}

		}

	}

	for _, schema := range ProgramsSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {

			_, ok1 := result_schemas["programs"]

			if ok1 == false {
				result_schemas["programs"] = make([]schemas.Tool, 0)
			}

			result_schemas["programs"] = append(result_schemas["programs"], schema)

			_, ok2 := result_tools["programs"]

			if ok2 == false {
				result_tools["programs"] = NewPrograms(playground, sandbox, allowed_programs)
			}

		}

	}

	for _, schema := range RequirementsSchema {

		if slices.Contains(allowed_tools, schema.Function.Name) {

			_, ok1 := result_schemas["requirements"]

			if ok1 == false {
				result_schemas["requirements"] = make([]schemas.Tool, 0)
			}

			result_schemas["requirements"] = append(result_schemas["requirements"], schema)

			_, ok2 := result_tools["requirements"]

			if ok2 == false {
				result_tools["requirements"] = NewRequirements(playground, sandbox)
			}

		}

	}

	return result_schemas, result_tools

}
