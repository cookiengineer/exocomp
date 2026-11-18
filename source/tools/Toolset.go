package tools

import "exocomp/types"
import utils_capabilities "exocomp/utils/capabilities"
import _ "embed"
import net_url "net/url"
import "strings"

func Toolset(
	playground string,
	sandbox string,
	model string,
	url *net_url.URL,
	debug bool,
	allowed_programs []string,
	allowed_images []string,
	allowed_tools []string,
	agent_name string,
	capabilities *utils_capabilities.Set,
) []types.Tool {

	filtered_by_methods := make(map[string][]string, 0)
	filtered_by_tool := make(map[string]types.Tool)

	for _, tool_name := range allowed_tools {

		if strings.Contains(tool_name, ".") {

			tmp1 := strings.TrimSpace(tool_name[0:strings.Index(tool_name, ".")])
			tmp2 := strings.TrimSpace(tool_name[strings.Index(tool_name, ".")+1:])

			name := strings.ToLower(tmp1)
			method := strings.ToUpper(tmp2[0:1]) + strings.ToLower(tmp2[1:])

			_, ok := filtered_by_methods[name]

			if ok == false {
				filtered_by_methods[name] = make([]string, 0)
			}

			filtered_by_methods[name] = append(filtered_by_methods[name], method)

		}

	}

	for tool_name, allowed_methods := range filtered_by_methods {

		var tool types.Tool = nil

		switch tool_name {

		case "agents":
			tool = NewAgents(allowed_methods, playground, sandbox, model, url, debug)
		case "bugs":
			tool = NewBugs(allowed_methods, playground, sandbox)
		case "changelog":
			tool = NewChangelog(allowed_methods, playground, sandbox)
		case "containers":
			tool = NewContainers(allowed_methods, playground, sandbox, allowed_images)
		case "files":
			tool = NewFilesWithCapabilities(allowed_methods, playground, sandbox, capabilities)
		case "findings":
			tool = NewFindings(allowed_methods, playground, sandbox, agent_name)
		case "humans":
			tool = NewHumansWithCapabilities(allowed_methods, playground, sandbox, agent_name, capabilities)
		case "programs":
			tool = NewProgramsWithCapabilities(allowed_methods, playground, sandbox, allowed_programs, capabilities)
		case "requirements":
			tool = NewRequirements(allowed_methods, playground, sandbox)
		case "shell":
			tool = NewShellWithCapabilities(allowed_methods, playground, sandbox, allowed_programs, capabilities)
		case "skills":
			tool = NewSkillsWithCapabilities(allowed_methods, playground, sandbox, allowed_programs, allowed_tools, capabilities)
		case "websites":
			tool = NewWebsitesWithCapabilities(allowed_methods, playground, sandbox, capabilities)

		}

		if tool != nil {
			filtered_by_tool[tool_name] = tool
		}

	}

	result := make([]types.Tool, 0)

	for _, tool := range filtered_by_tool {
		result = append(result, tool)
	}

	return result

}
