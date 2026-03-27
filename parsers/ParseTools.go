package parsers

import "exocomp/tools"
import "fmt"
import "strings"

func ParseTools(agent string, sandbox string, allowed_tools []string, allowed_programs []string, response_content string) []tools.Tool {

	result := make([]tools.Tool, 0)

	lines := strings.Split(strings.TrimSpace(response_content), "\n")

	for l := 0; l < len(lines); l++ {

		line := strings.TrimSpace(lines[l])
		line = strings.ReplaceAll(line, "\r", "")

		if strings.HasPrefix(line, "#!tool:") && strings.Contains(line, " ") {

			tmp := strings.TrimSpace(line[7:strings.Index(line, " ")])

			if strings.Contains(tmp, ".") {

				tool_name := strings.TrimSpace(strings.ToLower(tmp[0:strings.Index(tmp, ".")]))

				if isAllowedTool(allowed_tools, tool_name) == true {

					tool_body := []string{line}

					if strings.HasSuffix(line, "<<#!EOF") {

						for s := l; s < len(lines); s++ {

							seek := strings.TrimSpace(lines[s])

							if seek == "#!EOF" {
								tool_body = append(tool_body, lines[s])
								break
							} else {
								tool_body = append(tool_body, lines[s])
							}

						}

					}

					if tool_name == "files" {

						files_tool := tools.NewFiles(agent, sandbox, allowed_tools, allowed_programs)
						tool, parsed, err := files_tool.Parse(strings.Join(tool_body, "\n"))

						if err == nil {

							if parsed[1] > 1 {
								l += parsed[1]
							}

							result = append(result, tool)

						} else {

							result = append(result, tools.Tool(tools.NewStub(
								tool_name,
								err,
							)))

						}

					} else if tool_name == "notes" {

						notes_tool := tools.NewNotes(agent, sandbox, allowed_tools, allowed_programs)
						tool, parsed, err := notes_tool.Parse(strings.Join(tool_body, "\n"))

						if err == nil {

							if parsed[1] > 1 {
								l += parsed[1]
							}

							result = append(result, tool)

						} else {

							result = append(result, tools.Tool(tools.NewStub(
								tool_name,
								err,
							)))

						}

					} else if tool_name == "programs" {

						programs_tool := tools.NewNotes(agent, sandbox, allowed_tools, allowed_programs)
						tool, parsed, err := programs_tool.Parse(strings.Join(tool_body, "\n"))

						if err == nil {

							if parsed[1] > 1 {
								l += parsed[1]
							}

							result = append(result, tool)

						} else {

							result = append(result, tools.Tool(tools.NewStub(
								tool_name,
								err,
							)))

						}

					} else if tool_name == "tasks" {

						// TODO: Integrate Tasks Tool
						result = append(result, tools.Tool(tools.NewStub(
							tool_name,
							fmt.Errorf("#!tool:%s: Tool not implemented yet.", tool_name),
						)))

					} else {

						result = append(result, tools.Tool(tools.NewStub(
							tool_name,
							fmt.Errorf("#!tool:%s: Tool does not exist.", tool_name),
						)))

					}

				} else {

					result = append(result, tools.Tool(tools.NewStub(
						tool_name,
						fmt.Errorf("#!tool:%s: Tool does not exist.", tool_name),
					)))

				}

			} else {

				result = append(result, tools.Tool(tools.NewStub(
					"???",
					fmt.Errorf("#!tool:???: Invalid Tool Call line."),
				)))

			}

		}

	}

	return result

}
