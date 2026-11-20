package skills

import "fmt"

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
