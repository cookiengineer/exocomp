package tools

import "exocomp/types"
import utils_bytes "exocomp/utils/bytes"
import utils_fmt "exocomp/utils/fmt"
import "context"
import "errors"
import "fmt"
import "io/fs"
import "os"
import "os/exec"
import "path/filepath"
import "sort"
import "strings"
import "time"

type Skills struct {
	AllowedPrograms []string
	AllowedTools    []string
	Playground      string
	Sandbox         string
	contents        map[string]*types.Skill
	loaded_skills   map[string]*types.Skill
	processes       map[string]*os.Process
}

func NewSkills(playground string, sandbox string, allowed_programs []string, allowed_tools []string) *Skills {

	skills := &Skills{
		AllowedPrograms: allowed_programs,
		AllowedTools:    allowed_tools,
		Playground:      playground,
		Sandbox:         sandbox,
		contents:        make(map[string]*types.Skill),
		loaded_skills:   make(map[string]*types.Skill),
		processes:       make(map[string]*os.Process),
	}

	// NOTE: readSkills() allowed only at bootup time
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

	} else if method == "Execute" {

		name,     ok1 := arguments["name"].(string)
		script,   ok2 := arguments["script"].(string)
		raw_args, ok3 := arguments["arguments"].([]interface{})

		if ok1 == true && ok2 == true && ok3 == true {

			args := make([]string, len(raw_args))

			for a, value := range raw_args {

				tmp, ok := value.(string)

				if ok == true {
					args[a] = tmp
				}

			}

			return tool.Execute(utils_fmt.FormatSkillName(name), script, args)

		} else if ok1 == true && ok2 == true && ok3 == false {
			return "", fmt.Errorf("skills.%s: %s", method, "Invalid parameter \"arguments\" is not an array of strings.")
		} else if ok1 == true && ok2 == false && ok3 == true {
			return "", fmt.Errorf("skills.%s: %s", method, "Invalid parameter \"script\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true {
			return "", fmt.Errorf("skills.%s: %s", method, "Invalid parameter \"name\" is not a string.")
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

		missing_programs := make([]string, 0)
		missing_tools    := make([]string, 0)

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

func (tool *Skills) Execute(name string, script string, arguments []string) (string, error) {

	skill, ok1 := tool.loaded_skills[name]

	if ok1 == true {

		runtime, ok2 := skill.Scripts[script]

		if ok2 == true {

			found := false

			for _, program := range tool.AllowedPrograms {

				if program == runtime {
					found = true
					break
				}

			}

			if found == true {

				script_path, err1 := sanitizeSandboxPath(tool.Playground, filepath.Join(tool.Playground, "skills", skill.Name, "scripts", script))

				if err1 == nil {

					runtime_arguments := make([]string, 0)

					if runtime == "go" {
						runtime_arguments = append(runtime_arguments, "run")
						runtime_arguments = append(runtime_arguments, script_path)
					} else {
						runtime_arguments = append(runtime_arguments, script_path)
					}

					for a := 0; a < len(arguments); a++ {

						if strings.Contains(arguments[a], string(os.PathSeparator)) {

							resolved, err := sanitizeSandboxPath(tool.Sandbox, arguments[a])

							if err == nil {
								runtime_arguments = append(runtime_arguments, resolved)
							} else {
								return "", fmt.Errorf("skills.Execute: %s", err.Error())
							}

						} else {
							runtime_arguments = append(runtime_arguments, arguments[a])
						}

					}

					ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Minute)
					buffer      := utils_bytes.NewContextBuffer(16*1024*1024, cancel)

					defer cancel()

					go func() {

						ticker := time.NewTicker(10 * time.Second)

						BackgroundLoop:
						for {

							select {

							case <-ctx.Done():

								break BackgroundLoop

							case <-ticker.C:

								last_write := buffer.LastWrite()

								if time.Since(last_write) > 1 * time.Minute {

									cancel()
									break BackgroundLoop

								}

							}

						}

						ticker.Stop()

					}()


					cmd    := exec.CommandContext(ctx, runtime, runtime_arguments...)
					cmd.Dir = tool.Sandbox

					cmd.Stdin  = strings.NewReader("")
					cmd.Stdout = buffer
					cmd.Stderr = buffer

					err2   := cmd.Run()
					result := strings.Join([]string{
						fmt.Sprintf("skills.Execute: %s %s", runtime, strings.Join(runtime_arguments, " ")),
						buffer.String(),
					}, "\n")

					if ctx.Err() == context.Canceled && buffer.IsTruncated() {

						return result, fmt.Errorf("skills.Execute: Script output exceeded 16MB limit")

					} else if ctx.Err() == context.DeadlineExceeded {

						return result, fmt.Errorf("skills.Execute: Script timeout exceeded 10mins limit")

					} else if err2 == nil {

						return result, nil

					} else {

						if errors.Is(err2, fs.ErrPermission) {
							return "", fmt.Errorf("skills.Execute: Invalid script \"%s\": Permission denied.", script_path)
						} else if errors.Is(err2, fs.ErrNotExist) || strings.Contains(err2.Error(), "executable file not found") {
							return "", fmt.Errorf("skills.Execute: Invalid runtime \"%s\": Program doesn't exist.", runtime)
						} else {
							return result, fmt.Errorf("skills.Execute: Runtime \"%s\" execution error \"%s\".", runtime, err2.Error())
						}

					}

				} else {
					return "", fmt.Errorf("skills.Execute: %s", err1.Error())
				}

			} else {
				return "", fmt.Errorf("skills.Execute: Invalid runtime \"%s\": Attempt to execute unallowed program", runtime)
			}


		} else {
			return "", fmt.Errorf("skills.Execute: Script \"%s\" has no runtime!", script)
		}

	} else {
		return "", fmt.Errorf("skills.Execute: Skill \"%s\" isn't loaded!", name)
	}

}
