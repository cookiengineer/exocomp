package commands

import "fmt"
import "strings"

func Resolve(catalog *Catalog, input string) (*Command, error) {

	command, raw, err0 := parseRaw(input)

	if err0 == nil {

		tool, ok := catalog.Tool(command.Name)

		if ok == true {

			arguments := make(map[string]any)
			properties := tool.Function.Parameters.Properties

			for key, value := range raw {

				property, has_property := properties[key]

				if has_property == true {

					parsed, err1 := ParseValue(value, &property)

					if err1 == nil {
						arguments[key] = parsed
					} else {
						return nil, fmt.Errorf("commands.Resolve: %s", err1.Error())
					}

				} else {
					return nil, fmt.Errorf("commands.Resolve: Unknown parameter \"%s\" for \"%s\".", key, tool.Function.Name)
				}

			}

			for _, required := range tool.Function.Parameters.Required {

				_, has_value := arguments[required]

				if has_value == false {
					return nil, fmt.Errorf("commands.Resolve: Missing required parameter \"%s\" for \"%s\".", required, tool.Function.Name)
				}

			}

			name := tool.Function.Name
			method := name[strings.LastIndex(name, ".")+1:]

			return &Command{
				Name:      name,
				Method:    method,
				Arguments: arguments,
			}, nil

		} else {
			return nil, fmt.Errorf("commands.Resolve: Unknown tool \"%s\".", command.Name)
		}

	} else {
		return nil, err0
	}

}
