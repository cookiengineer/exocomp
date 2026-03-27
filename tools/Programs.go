package tools

import "exocomp/utils"
import "bytes"
import "fmt"
import "os"
import "os/exec"
import "slices"
import "strings"

type Programs struct {
	Method    string
	Arguments []string
	Programs  []string
	Sandbox   string
}

func NewPrograms(agent string, sandbox string, tools []string, programs []string) *Programs {

	return &Programs{
		Method:    "",
		Arguments: make([]string, 0),
		Programs:  programs,
		Sandbox:   sandbox,
	}

}

func (tool *Programs) Call() (string, error) {

	if tool.Method == "Help" {
		return tool.Help(tool.Arguments)
	} else if tool.Method == "List" {
		return tool.List(tool.Arguments)
	} else if tool.Method == "Execute" {
		return tool.Execute(tool.Arguments)
	} else {
		return "", fmt.Errorf("#!tool:programs.%s: Invalid method.", tool.Method)
	}

}

func (tool *Programs) Help(arguments []string) (string, error) {

	return strings.Join([]string{
		"#!tool:programs.List",
		"",
		"#!tool:programs.Execute \"program-name\" arg1 \"arg2 with whitespace\" arg3",
	}, "\n"), nil

}

func (tool *Programs) Execute(arguments []string) (string, error) {

	if len(arguments) > 1 {

		if slices.Contains(tool.Programs, arguments[0]) == true {

			program           := arguments[0]
			program_arguments := make([]string, 0)

			for a := 1; a < len(arguments); a++ {

				if strings.Contains(arguments[a], string(os.PathSeparator)) {

					resolved, err := resolveSandboxPath(tool.Sandbox, arguments[a])

					if err == nil {
						program_arguments = append(program_arguments, resolved)
					} else {
						return "", fmt.Errorf("#!tool:programs.Execute: %s", err.Error())
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
					fmt.Sprintf("#!tool:programs.Execute: %s %s", program, strings.Join(program_arguments, " ")),
					buffer.String(),
				}, "\n")

				return result, nil

			} else {

				result := strings.Join([]string{
					fmt.Sprintf("#!tool:programs.Execute: %s %s", program, strings.Join(program_arguments, " ")),
					err.Error(),
				}, "\n")

				return result, nil

			}

		} else {
			return "", fmt.Errorf("#!tool:programs.Execute: Invalid program \"%s\", must be an allowed program.", arguments[0])
		}

	} else {
		return "", fmt.Errorf("#!tool:programs.Execute: Invalid arguments, at least one argument required.")
	}

}

func (tool *Programs) List(arguments []string) (string, error) {

	if len(arguments) == 0 {

		names := make([]string, 0)

		for _, name := range tool.Programs {
			names = append(names, name)
		}

		slices.Sort(names)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("#!tool:programs.List:"))

		for _, name := range names {
			result = append(result, fmt.Sprintf("Name: %s", name))
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("#!tool:programs.List: Only zero arguments allowed")
	}

}

func (tool *Programs) Parse(text string) (Tool, [2]int, error) {

	// #!tool:programs.List
	// #!tool:programs.Execute arg1 "arg2 with whitespace" arg3...

	lines := strings.Split(text, "\n")

	if len(lines) > 0 && strings.HasPrefix(lines[0], "#!tool:programs.") {

		fields := utils.SplitArguments(strings.TrimSpace(lines[0][len("#!tool:programs."):]))
		method := strings.ToUpper(fields[0][0:1]) + strings.ToLower(fields[0][1:])
		parsed := [2]int{0, 1}

		if method == "Help" ||
			method == "List" ||
			method == "Execute" {

			tool.Method    = method
			tool.Arguments = fields[1:]

			return Tool(tool), parsed, nil

		} else {
			return nil, [2]int{0, len(lines)}, fmt.Errorf("#!tool:programs.%s: Invalid method.", method)
		}

	} else {
		return nil, [2]int{0, len(lines)}, fmt.Errorf("Invalid Tool Call line.")
	}

}
