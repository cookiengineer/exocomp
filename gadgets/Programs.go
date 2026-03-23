package gadgets

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

func NewPrograms(sandbox string, programs []string) *Programs {

	return &Programs{
		Programs: programs,
		Sandbox:  sandbox,
	}

}

func (gadget *Programs) Help(arguments []string) (string, error) {

	return strings.Join([]string{
		"#!gadget:programs.List",
		"",
		"#!gadget:programs.Execute \"program-name\" \"arg1\" \"arg2\" \"arg3\"",
	}, "\n"), nil

}

func (gadget *Programs) List(arguments []string) (string, error) {

	if len(arguments) == 0 {

		prompt := "The list of available programs is:"

		for _, name := range gadget.Programs {
			prompt += name + "\n"
		}

		return prompt, nil

	} else {
		return "", fmt.Errorf("#!programs.List: Only zero arguments allowed")
	}

}

func (gadget *Programs) Execute(arguments []string) (string, error) {

	if len(arguments) > 1 && slices.Contains(gadget.Programs, arguments[0]) == true {

		program           := arguments[0]
		program_arguments := make([]string, 0)

		for a := 1; a < len(arguments); a++ {

			if strings.Contains(arguments[a], string(os.PathSeparator)) {

				resolved, err := resolveSandboxPath(gadget.Sandbox, arguments[a])

				if err == nil {
					program_arguments = append(program_arguments, resolved)
				} else {
					return "", fmt.Errorf("#!programs.Execute: %s", err.Error())
				}

			}

		}

		buffer    := bytes.Buffer{}
		cmd       := exec.Command(program, program_arguments...)
		cmd.Dir    = gadget.Sandbox
		cmd.Stdout = &buffer
		cmd.Stderr = &buffer

		err := cmd.Run()

		if err == nil {

			result := strings.Join([]string{
				fmt.Sprintf("#!programs.Execute: %s %s", program, strings.Join(program_arguments, " ")),
				buffer.String(),
			}, "\n")

			return result, nil

		} else {

			result := strings.Join([]string{
				fmt.Sprintf("#!programs.Execute: %s %s", program, strings.Join(program_arguments, " ")),
				err.Error(),
			}, "\n")

			return result, nil

		}

	} else {
		return "", fmt.Errorf("#!programs.Execute: %s must be an allowed program", arguments[0])
	}

}
