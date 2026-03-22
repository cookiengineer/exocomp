package gadgets

import "exocomp/config"
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

func NewPrograms(config *config.Config) *Programs {

	return &Programs{
		Programs: config.Programs,
		Sandbox:  config.Sandbox,
	}

}

func (gadget *Programs) Help() string {

	return strings.Join([]string{
		"Programs Gadget Usage:",
		"",
		"List programs:",
		"#!gadget:programs.List",
		"",
		"Execute programs ",
		"#!gadget:programs.Execute \"program-name\" \"arg1\" \"arg2\" \"arg3\"",
		"",
	}, "\n")

}

func (gadget *Programs) List() (string, error) {

	// TODO: List available programs

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
