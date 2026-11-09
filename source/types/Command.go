package types

import "fmt"
import "strings"

type Command struct {
	Name      string         `json:"name"`
	Method    string         `json:"method"`
	Arguments map[string]any `json:"arguments"`
}

func ParseCommand(prompt string) *Command {

	command := &Command{
		Name:      "",
		Method:    "",
		Arguments: make(map[string]any),
	}

	err := command.Parse(prompt)

	if err == nil {
		return command
	}

	return nil

}

func (command *Command) Parse(prompt string) error {

	if strings.HasPrefix(prompt, "/") && strings.Contains(prompt, " ") && !strings.Contains(prompt, "\n") {

		name := prompt[1:strings.Index(prompt, " ")]

		if strings.Contains(name, ".") {

			command.Name = strings.TrimSpace(name)
			command.Method = strings.TrimSpace(command.Name[strings.LastIndex(command.Name, ".")+1:])
			command.Arguments = parseCommandArguments(strings.TrimSpace(prompt[1+len(name)+1:]))

			return nil

		} else {
			return fmt.Errorf("types.ParseCommand: Invalid command tool name \"%s\".", name)
		}

	} else {
		return fmt.Errorf("types.ParseCommand: Invalid command prompt syntax.")
	}

}
