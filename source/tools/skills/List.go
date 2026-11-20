package skills

import "fmt"
import "sort"
import "strings"

func (tool *Skills) List() (string, error) {

	readSkills(tool)

	if len(tool.contents) > 0 {

		lines := make([]string, 0)

		for name, skill := range tool.contents {

			status := "unloaded"
			scripts := make([]string, 0)
			tools := make([]string, 0)
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
