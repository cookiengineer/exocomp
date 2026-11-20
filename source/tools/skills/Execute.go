package skills

import utils_bytes "exocomp/utils/bytes"
import utils_sandbox "exocomp/utils/sandbox"
import "context"
import "fmt"
import "os"
import "path/filepath"
import "strings"
import "time"

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

				script_path, err1 := utils_sandbox.SanitizePath(tool.Playground, filepath.Join(tool.Playground, "skills", skill.Name, "scripts", script))

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

							resolved, err := utils_sandbox.SanitizePath(tool.Sandbox, arguments[a])

							if err == nil {
								runtime_arguments = append(runtime_arguments, resolved)
							} else {
								return "", fmt.Errorf("skills.Execute: %s", err.Error())
							}

						} else {
							runtime_arguments = append(runtime_arguments, arguments[a])
						}

					}

					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
					buffer := utils_bytes.NewContextBuffer(16*1024*1024, cancel)

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

								if time.Since(last_write) > 1*time.Minute {

									cancel()
									break BackgroundLoop

								}

							}

						}

						ticker.Stop()

					}()

					scripts_dir := filepath.Join(tool.Playground, "skills", skill.Name, "scripts")
					read_paths, write_paths := tool.Capabilities.SandboxPaths()
					read_paths = append(read_paths, scripts_dir)

					cmd, err2 := utils_sandbox.CommandContextWithPermissions(ctx, tool.Sandbox, runtime, runtime_arguments, utils_sandbox.Permissions{
						ReadPaths:  read_paths,
						WritePaths: write_paths,
						Network:    tool.Capabilities.AllowsNetworkEgress(),
					})

					if err2 == nil {

						cmd.Stdin = strings.NewReader("")
						cmd.Stdout = buffer
						cmd.Stderr = buffer

						err3 := cmd.Run()
						result := strings.Join([]string{
							fmt.Sprintf("skills.Execute: %s %s", runtime, strings.Join(runtime_arguments, " ")),
							buffer.String(),
						}, "\n")

						if ctx.Err() == context.Canceled && buffer.IsTruncated() {

							return result, fmt.Errorf("skills.Execute: Script output exceeded 16MB limit")

						} else if ctx.Err() == context.DeadlineExceeded {

							return result, fmt.Errorf("skills.Execute: Script timeout exceeded 10mins limit")

						} else if err3 == nil {

							return result, nil

						} else {
							return utils_sandbox.SanitizeExecutionError("skills", "Execute", "runtime", runtime, result, err3)
						}

					} else {
						return "", fmt.Errorf("skills.Execute: %s", err2.Error())
					}

				} else {
					return "", fmt.Errorf("skills.Execute: Invalid script \"%s\".", script)
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
