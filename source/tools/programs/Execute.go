package programs

import utils_bytes "exocomp/utils/bytes"
import utils_sandbox "exocomp/utils/sandbox"
import utils_shell "exocomp/utils/shell"
import "context"
import "fmt"
import "slices"
import "strings"
import "time"

func (tool *Programs) Execute(program string, arguments []string) (string, error) {

	if slices.Contains(tool.AllowedPrograms, program) {

		if strings.ContainsAny(program, "/\\") == false {

			program_arguments := make([]string, 0)

			for a := 0; a < len(arguments); a++ {

				argument, err := utils_shell.SanitizeProgramArgument(tool.Sandbox, tool.Capabilities, arguments[a])

				if err == nil {
					program_arguments = append(program_arguments, argument)
				} else {
					return "", fmt.Errorf("programs.Execute: %s", err.Error())
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

			read_paths, write_paths := tool.Capabilities.SandboxPaths()

			cmd, err2 := utils_sandbox.CommandContextWithPermissions(ctx, tool.Sandbox, program, program_arguments, utils_sandbox.Permissions{
				ReadPaths:  read_paths,
				WritePaths: write_paths,
				Network:    tool.Capabilities.AllowsNetworkEgress(),
			})

			if err2 == nil {

				cmd.Stdin = strings.NewReader("")
				cmd.Stdout = buffer
				cmd.Stderr = buffer

				result := ""
				err3 := cmd.Run()

				if len(program_arguments) > 0 {
					result = strings.Join([]string{
						fmt.Sprintf("programs.Execute: %s %s", program, strings.Join(program_arguments, " ")),
						buffer.String(),
					}, "\n")
				} else {
					result = strings.Join([]string{
						fmt.Sprintf("programs.Execute: %s", program),
						buffer.String(),
					}, "\n")
				}

				if ctx.Err() == context.Canceled && buffer.IsTruncated() {

					return result, fmt.Errorf("programs.Execute: Program output exceeded 16MB limit")

				} else if ctx.Err() == context.DeadlineExceeded {

					return result, fmt.Errorf("programs.Execute: Program timeout exceeded 10mins limit")

				} else if err3 == nil {

					return result, nil

				} else {
					return utils_sandbox.SanitizeExecutionError("programs", "Execute", "program", program, result, err3)
				}

			} else {
				return "", fmt.Errorf("programs.Execute: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("programs.Execute: Invalid program \"%s\": Attempt to execute unallowed program", program)
		}
	} else {
		return "", fmt.Errorf("programs.Execute: Invalid program \"%s\": Attempt to execute unallowed program", program)
	}

}
