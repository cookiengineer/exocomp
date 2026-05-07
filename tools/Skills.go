package tools

import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"
import "fmt"
import "os"

type Skills struct {
	Playground    string
	Sandbox       string
	Skills        map[string]*types.Skill
	loaded_skills map[string]*types.Skill
	processes     map[string]*os.Process
}

func NewSkills(playground string, sandbox string) *Skills {

	skills := &Skills{
		Playground:    playground,
		Sandbox:       sandbox,
		Skills:        make(map[string]*types.Skill),
		loaded_skills: make(map[string]*types.Skill),
		processes:     make(map[string]*os.Process),
	}

	// TODO: Read folder at bootup time or at runtime?

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

func (tool *Skills) List() (string, error) {

	// TODO: Read playground/skills folder
	// TODO: Parse each SKILL.md's header --- for yaml metadata
	// TODO: List skills in an overview, Name, Description, Compatibility, Allowed-Tools
	// TODO: List also Scripts as whatever.py, foo.js etc

	return "", fmt.Errorf("skills.Load: %s", "Not implemented")

}

func (tool *Skills) Load(name string) (string, error) {

	// TODO: Read the SKILL.md body
	// and return body markdown as result, nil

	return "", fmt.Errorf("skills.Load: %s", "Not implemented")

}

func (tool *Skills) ExecuteScript(script string, arguments []string) (string, error) {

	// TODO: Execute script of the skill

	return "", fmt.Errorf("skills.Load: %s", "Not implemented")

}
