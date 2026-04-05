package tools

import "bytes"
import "fmt"
import "os"
import "os/exec"
import "slices"
import "strings"

type Programs struct {
	Programs []string
	Sandbox  string
}

func NewPrograms(agent string, sandbox string, programs []string) *Programs {

	return &Programs{
		Programs: programs,
		Sandbox:  sandbox,
	}

}

func (tool *Programs) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "List" {

		return tool.List()

	} else if method == "Execute" {

		program,  ok1 := arguments["program"].(string)
		raw_args, ok2 := arguments["arguments"].([]interface{})

		if ok1 == true && ok2 == true {

			args := make([]string, len(raw_args))

			for a, value := range raw_args {

				tmp, ok := value.(string)

				if ok == true {
					args[a] = tmp
				}

			}

			return tool.Execute(program, args)

		} else {
			return "", fmt.Errorf("programs.Execute: Invalid parameters")
		}

	} else {
		return "", fmt.Errorf("programs.%s: Invalid method.", method)
	}

}

func (tool *Programs) Execute(program string, arguments []string) (string, error) {

	if slices.Contains(tool.Programs, program) {

		program_arguments := make([]string, 0)

		for a := 0; a < len(arguments); a++ {

			if strings.Contains(arguments[a], string(os.PathSeparator)) {

				resolved, err := sanitizeSandboxPath(tool.Sandbox, arguments[a])

				if err == nil {
					program_arguments = append(program_arguments, resolved)
				} else {
					return "", fmt.Errorf("programs.Execute: %s", err.Error())
				}

			}

		}

		buffer    := bytes.Buffer{}
		cmd       := exec.Command(program, program_arguments...)
		cmd.Dir    = tool.Sandbox
		cmd.Stdout = &buffer
		cmd.Stderr = &buffer

		err := cmd.Run()

		if err == nil {

			result := strings.Join([]string{
				fmt.Sprintf("programs.Execute: %s %s", program, strings.Join(program_arguments, " ")),
				buffer.String(),
			}, "\n")

			return result, nil

		} else {

			result := strings.Join([]string{
				fmt.Sprintf("programs.Execute: %s %s", program, strings.Join(program_arguments, " ")),
				err.Error(),
			}, "\n")

			return result, nil

		}

	} else {
		return "", fmt.Errorf("programs.Execute: Invalid program \"%s\", must be an allowed program.", program)
	}

}

func (tool *Programs) List() (string, error) {

	names := make([]string, 0)

	for _, name := range tool.Programs {
		names = append(names, name)
	}

	slices.Sort(names)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("programs.List:"))

	for _, name := range names {
		result = append(result, fmt.Sprintf("Name: %s", name))
	}

	return strings.Join(result, "\n"), nil

}

