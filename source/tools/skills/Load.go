package skills

import "fmt"
import "strings"

func (tool *Skills) Load(name string) (string, error) {

	skill, ok := tool.contents[name]

	if ok == true {

		missing_programs := make([]string, 0)
		missing_tools := make([]string, 0)

		if len(skill.AllowedPrograms) > 0 {

			for _, program_name := range skill.AllowedPrograms {

				found := false

				for _, program := range tool.AllowedPrograms {

					if program == program_name {

						found = true
						break

					}

				}

				if found == false {
					missing_programs = append(missing_programs, program_name)
				}

			}

		}

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

		if len(missing_tools) == 0 && len(missing_programs) == 0 {

			tool.loaded_skills[skill.Name] = skill

			// NOTE: Session.LoadSkill() does actual loading
			return fmt.Sprintf("skills.Load: Skill \"%s\" got loaded.", name), nil

		} else if len(missing_programs) != 0 {
			return "", fmt.Errorf("skills.Load: Can't load Skill because of missing Programs %s", strings.Join(missing_programs, " and "))
		} else if len(missing_tools) != 0 {
			return "", fmt.Errorf("skills.Load: Can't load Skill because of missing Tools %s", strings.Join(missing_tools, " and "))
		} else {
			return "", fmt.Errorf("skills.Load: Can't load Skill \"%s\"", name)
		}

	} else {
		return "", fmt.Errorf("skills.Load: Skill \"%s\" doesn't exist!", name)
	}

}
