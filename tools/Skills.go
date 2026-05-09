package tools

import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"
import "fmt"
import "os"
import "sort"
import "strings"

type Skills struct {
	AllowedTools  []string
	Playground    string
	Sandbox       string
	contents      map[string]*types.Skill
	loaded_skills map[string]*types.Skill
	processes     map[string]*os.Process
}

func NewSkills(playground string, sandbox string, allowed_tools []string) *Skills {

	skills := &Skills{
		AllowedTools:  allowed_tools,
		Playground:    playground,
		Sandbox:       sandbox,
		contents:      make(map[string]*types.Skill),
		loaded_skills: make(map[string]*types.Skill),
		processes:     make(map[string]*os.Process),
	}

	// NOTE: readSkills() only at bootup time in case Agent rewrites SKILL.md at runtime
	readSkills(skills)

	return skills

}

func (tool *Skills) Call(method string, arguments map[string]interface{}) (string, error) {

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

	} else if method == "ExecuteScript" {

		script,   ok1 := arguments["script"].(string)
		raw_args, ok2 := arguments["arguments"].([]interface{})

		if ok1 == true && ok2 == true {

			args := make([]string, len(raw_args))

			for a, value := range raw_args {

				tmp, ok := value.(string)

				if ok == true {
					args[a] = tmp
				}

			}

			return tool.ExecuteScript(script, args)

		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("skills.%s: %s", method, "Invalid parameter \"arguments\" is not an array of strings.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("skills.%s: %s", method, "Invalid parameter \"script\" is not a string.")
		} else {
			return "", fmt.Errorf("skills.%s: Invalid parameters.", method)
		}

	} else {
		return "", fmt.Errorf("skills.%s: Invalid method.", method)
	}

}

func (tool *Skills) Get(id string) (any, error) {

	name        := utils_fmt.FormatSkillName(id)
	content, ok := tool.contents[name]

	if ok == true {
		return content, nil
	} else {
		return nil, fmt.Errorf("skills.Get: No skill found with the name \"%s\".", name)
	}

}

func (tool *Skills) List() (string, error) {

	readSkills(tool)

	if len(tool.contents) > 0 {

		lines := make([]string, 0)

		for name, skill := range tool.contents {

			status  := "unloaded"
			scripts := make([]string, 0)
			tools   := make([]string, 0)
			tmp, ok := tool.loaded_skills[name]

			if ok == true && tmp != nil {
				status = "loaded"
			}

			if len(skill.AllowedTools) > 0 {

				for _, tool := range skill.AllowedTools {
					tools = append(tools, tool)
				}

				sort.Strings(tools)

			}

			if len(skill.Scripts) > 0 {

				for script, _ := range skill.Scripts {
					scripts = append(scripts, script)
				}

				sort.Strings(scripts)

			}

			lines = append(lines, fmt.Sprintf("- Skill: %s, Status: %s, Description: %s, Tools: %s, Scripts: %s", skill.Name, status, skill.Description, strings.Join(tools, " "), strings.Join(scripts, " ")))

		}

		sort.Strings(lines)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("skills.List: %d skills available.", len(lines)))

		for l := 0; l < len(lines); l++ {
			result = append(result, lines[l])
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("skills.List: No skills available!")
	}

}

func (tool *Skills) Load(name string) (string, error) {

	skill, ok := tool.contents[name]

	if ok == true {

		missing_tools := make([]string, 0)

		if len(skill.AllowedTools) > 0 {

			for _, tool_name := range skill.AllowedTools {

				found := false

				for _, tool := range tool.AllowedTools {

					if tool == tool_name {
						found = true
						break
					}

				}

				if found == false {
					missing_tools = append(missing_tools, tool_name)
				}

			}

		}

		if len(missing_tools) == 0 {

			tool.loaded_skills[skill.Name] = skill

			// NOTE: Session.LoadSkill() does actual loading
			return fmt.Sprintf("skills.Load: Skill \"%s\" got loaded.", name), nil

		} else {
			return "", fmt.Errorf("skills.Load: Can't load Skill because of missing Tools %s", strings.Join(missing_tools, " and "))
		}

	} else {
		return "", fmt.Errorf("skills.Load: Skill \"%s\" doesn't exist!", name)
	}

}

func (tool *Skills) Unload(name string) (string, error) {

	skill, ok := tool.loaded_skills[name]

	if ok == true {

		delete(tool.loaded_skills, skill.Name)

		// NOTE: Session.UnloadSkill() does actual unloading
		return fmt.Sprintf("skills.Load: Skill \"%s\" got unloaded.", name), nil

	} else {
		return "", fmt.Errorf("skills.Load: Skill \"%s\" isn't loaded!", name)
	}

}

func (tool *Skills) ExecuteScript(script string, arguments []string) (string, error) {

	// TODO: Execute script of the skill

	return "", fmt.Errorf("skills.Load: %s", "Not implemented")

}
