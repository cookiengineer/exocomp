package tools

import utils_skill "exocomp/utils/skill"
import "fmt"
import "os"
import "path/filepath"
import "strings"

func readSkills(tool *Skills) error {

	if tool.Playground != "" {

		resolved, err0 := resolveSandboxPath(tool.Playground, "skills")

		if err0 == nil {

			for name, _ := range tool.Skills {
				delete(tool.Skills, name)
			}

			skills_entries, err1 := os.ReadDir(resolved)

			if err1 == nil {

				errors := make([]error, 0)

				for _, skill_entry := range skills_entries {

					if skill_entry.IsDir() {

						skill_name       := skill_entry.Name()
						skill_path       := filepath.Join(tool.Playground, "skills", skill_name, "SKILL.md")
						skill_stat, err1 := os.Stat(skill_path)

						if err1 == nil && skill_stat.IsDir() == false {

							skill_bytes, err2 := os.ReadFile(skill_path)

							if err2 == nil {

								skill := utils_skill.ParseSkill(skill_bytes)

								if skill != nil {

									scripts_path       := filepath.Join(tool.Playground, "skills", skill_name, "scripts")
									scripts_stat, err3 := os.Stat(scripts_path)

									if err3 == nil && scripts_stat.IsDir() {

										scripts_entries, err4 := os.ReadDir(scripts_path)

										if err4 == nil {

											for _, script_entry := range scripts_entries {

												if script_entry.IsDir() == false {

													script_name := script_entry.Name()

													if strings.HasSuffix(script_name, ".go") {
														skill.Scripts[script_name] = "go"
													} else {
														errors = append(errors, fmt.Errorf("Invalid Skill %s: Unsupported script \"%s\"", skill_name, script_name))
													}

												}

											}

										}

									}

									if skill.Name == skill_name {
										tool.Skills[skill.Name] = skill
									}

								} else {
									errors = append(errors, fmt.Errorf("Invalid Skill %s: %s", skill_name, "Cannot parse SKILL.md"))
								}

							}

						}

					}

				}

				if len(errors) > 0 {

					error_messages := make([]string, 0)

					for _, err := range errors {
						error_messages = append(error_messages, strings.TrimSpace(err.Error()))
					}

					return fmt.Errorf("readSkills: %s", strings.Join(error_messages, " | "))

				} else {
					return nil
				}

			} else {
				return fmt.Errorf("readSkills: %s", err1.Error())
			}

		} else {
			return fmt.Errorf("readSkills: %s", err0.Error())
		}

	} else {
		return fmt.Errorf("readSkills: Invalid tool playground")
	}

}
