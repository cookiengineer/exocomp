package commands

import "fmt"
import "strings"
import "unicode"

func parseRaw(input string) (*Command, map[string]string, error) {

	trimmed := strings.TrimSpace(input)

	if strings.HasPrefix(trimmed, "/") == false {
		return nil, nil, fmt.Errorf("commands.Parse: Command must start with \"/\".")
	}

	value := strings.TrimSpace(trimmed[1:])

	if value == "" {
		return nil, nil, fmt.Errorf("commands.Parse: Command name is empty.")
	}

	name := value
	arguments := ""

	for index, character := range value {

		if unicode.IsSpace(character) == true {

			name = value[0:index]
			arguments = strings.TrimSpace(value[index:])
			break

		}

	}

	if strings.Contains(name, ".") == false {
		return nil, nil, fmt.Errorf("commands.Parse: Invalid tool name \"%s\".", name)
	}

	method := name[strings.LastIndex(name, ".")+1:]

	raw, err := ParseRawArguments(arguments)

	if err == nil {

		return &Command{
			Name:      name,
			Method:    method,
			Arguments: make(map[string]any),
		}, raw, nil

	} else {
		return nil, nil, err
	}

}
