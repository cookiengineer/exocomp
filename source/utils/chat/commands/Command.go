package commands

type Command struct {
	Name      string
	Method    string
	Arguments map[string]any
}

func Parse(input string) (*Command, error) {

	command, raw, err := parseRaw(input)

	if err == nil {

		for key, value := range raw {

			parsed, err0 := ParseValue(value, nil)

			if err0 == nil {
				command.Arguments[key] = parsed
			} else {
				return nil, err0
			}

		}

		return command, nil

	} else {
		return nil, err
	}

}
