package skills

import "exocomp/schemas"
import "exocomp/types"
import utils_capabilities "exocomp/utils/capabilities"
import utils_fmt "exocomp/utils/fmt"
import "fmt"
import "os"
import "slices"
import "sort"

type Skills struct {
	Methods         []string
	Playground      string
	Sandbox         string
	AllowedPrograms []string
	AllowedTools    []string
	Capabilities    *utils_capabilities.Set
	contents        map[string]*types.Skill
	loaded_skills   map[string]*types.Skill
	processes       map[string]*os.Process
}

func NewSkills(methods []string, playground string, sandbox string, allowed_programs []string, allowed_tools []string) *Skills {

	return NewSkillsWithCapabilities(methods, playground, sandbox, allowed_programs, allowed_tools, nil)

}

func NewSkillsWithCapabilities(
	methods []string,
	playground string,
	sandbox string,
	allowed_programs []string,
	allowed_tools []string,
	capabilities *utils_capabilities.Set,
) *Skills {

	skills := &Skills{
		Methods:         methods,
		Playground:      playground,
		Sandbox:         sandbox,
		AllowedPrograms: allowed_programs,
		AllowedTools:    allowed_tools,
		Capabilities:    capabilities,
		contents:        make(map[string]*types.Skill),
		loaded_skills:   make(map[string]*types.Skill),
		processes:       make(map[string]*os.Process),
	}

	// NOTE: readSkills() allowed only at bootup time
	readSkills(skills)

	return skills

}

func (tool *Skills) Name() string {
	return "skills"
}

func (tool *Skills) Call(method string, arguments map[string]interface{}) (string, error) {

	if tool.HasMethod(method) == true {

		if method == "List" {

			return tool.List()

		} else if method == "Load" {

			name, ok1 := arguments["name"].(string)

			if ok1 == true {
				return tool.Load(utils_fmt.FormatSkillName(name))
			} else {
				return "", fmt.Errorf("skills.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			}

		} else if method == "Unload" {

			name, ok1 := arguments["name"].(string)

			if ok1 == true {
				return tool.Unload(utils_fmt.FormatSkillName(name))
			} else {
				return "", fmt.Errorf("skills.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			}

		} else if method == "Execute" {

			name, ok1 := arguments["name"].(string)
			script, ok2 := arguments["script"].(string)

			if ok1 == true && ok2 == true {

				args := make([]string, 0)

				raw, ok3 := arguments["arguments"]

				if ok3 == true {

					raw_args, ok4 := raw.([]interface{})

					if ok4 == true {

						args = make([]string, len(raw_args))

						for a, value := range raw_args {

							tmp, ok := value.(string)

							if ok == true {
								args[a] = tmp
							}

						}

					} else {
						return "", fmt.Errorf("skills.%s: %s", method, "Invalid parameter \"arguments\" is not an array of strings.")
					}

				}

				return tool.Execute(utils_fmt.FormatSkillName(name), script, args)

			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("skills.%s: %s", method, "Invalid parameter \"script\" is not a string.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("skills.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			} else {
				return "", fmt.Errorf("skills.%s: Invalid parameters.", method)
			}

		} else {
			return "", fmt.Errorf("skills.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("skills.%s: Method not allowed.", method)
	}

}

func (tool *Skills) GetContent(id string) (any, error) {

	name := utils_fmt.FormatSkillName(id)
	content, ok := tool.contents[name]

	if ok == true {
		return content, nil
	} else {
		return nil, fmt.Errorf("skills.Get: No skill found with the name \"%s\".", name)
	}

}

func (tool *Skills) GetContentIdentifiers() []string {

	result := make([]string, 0)

	for id, _ := range tool.contents {
		result = append(result, id)
	}

	sort.Strings(result)

	return result

}

func (tool *Skills) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Skills) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range SkillsSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}
