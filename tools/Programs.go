package tools

import "bytes"
import "errors"
import "fmt"
import "io/fs"
import "os"
import "os/exec"
import "slices"
import "strings"

type Programs struct {
	Programs []string
	Sandbox  string
}

func NewPrograms(playground string, sandbox string, allowed_programs []string) *Programs {

	return &Programs{
		Programs: allowed_programs,
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

		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("programs.%s: %s", method, "Invalid parameter \"arguments\" is not an array of strings.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("programs.%s: %s", method, "Invalid parameter \"program\" is not a string.")
		} else {
			return "", fmt.Errorf("programs.%s: Invalid parameters.", method)
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

			} else {
				program_arguments = append(program_arguments, arguments[a])
			}

		}

		buffer    := bytes.Buffer{}
		cmd       := exec.Command(program, program_arguments...)
		cmd.Dir    = tool.Sandbox
		cmd.Stdout = &buffer
		cmd.Stderr = &buffer

		err := cmd.Run()

		// TODO: Better errors for permission denied
		// fmt.Println("RUN ERROR", err)
		// fmt.Println(program, program_arguments)
		// fmt.Println(cmd.Dir)

		first_line := ""

		if len(program_arguments) > 0 {
			first_line = fmt.Sprintf("programs.Execute: %s %s", program, strings.Join(program_arguments, " "))
		} else {
			first_line = fmt.Sprintf("programs.Execute: %s", program)
		}

		if err == nil {

			result := strings.Join([]string{
				first_line,
				buffer.String(),
			}, "\n")

			return result, nil

		} else {

			if errors.Is(err, fs.ErrPermission) {
				return "", fmt.Errorf("programs.Execute: Invalid program \"%s\": Permission denied.", program)
			} else if errors.Is(err, fs.ErrNotExist) || strings.Contains(err.Error(), "executable file not found") {
				return "", fmt.Errorf("programs.Execute: Invalid program \"%s\": Program doesn't exist.", program)
			} else {

				result := strings.Join([]string{
					first_line,
					buffer.String(),
				}, "\n")

				return result, fmt.Errorf("programs.Execute: Program \"%s\" execution error \"%s\".", program, err.Error())

			}

		}

	} else {
		return "", fmt.Errorf("programs.Execute: Invalid program \"%s\": Attempt to execute unallowed program", program)
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

