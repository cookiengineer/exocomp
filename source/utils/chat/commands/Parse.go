package commands

import "fmt"
import "strings"
import "unicode"

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

func ParseRawArguments(raw string) (map[string]string, error) {

	result := make(map[string]string)
	index := 0

	for index < len(raw) {

		skipWhitespace(raw, &index)

		if index >= len(raw) {
			break
		}

		key_start := index

		for index < len(raw) && raw[index] != '=' && unicode.IsSpace(rune(raw[index])) == false {
			index++
		}

		if index >= len(raw) || raw[index] != '=' {
			return nil, fmt.Errorf("commands.Parse: Expected \"=\" after parameter \"%s\".", raw[key_start:index])
		}

		key := strings.TrimSpace(raw[key_start:index])
		index++

		if key == "" {
			return nil, fmt.Errorf("commands.Parse: Parameter name is empty.")
		}

		skipWhitespace(raw, &index)

		value, err := scanValue(raw, &index)

		if err == nil {
			result[key] = value
		} else {
			return nil, err
		}

	}

	return result, nil

}

func scanValue(raw string, index *int) (string, error) {

	if *index >= len(raw) {
		return "", nil
	}

	start := *index
	character := raw[*index]

	if character == '"' || character == '\'' {

		token := character
		*index++

		escaped := false

		for *index < len(raw) {

			current := raw[*index]

			if current == '\\' && escaped == false {
				escaped = true
				*index++
				continue
			}

			if current == token && escaped == false {
				*index++
				return raw[start:*index], nil
			}

			escaped = false
			*index++

		}

		return "", fmt.Errorf("commands.Parse: Unterminated string.")

	}

	if character == '{' || character == '[' {

		close_token := byte('}')

		if character == '[' {
			close_token = ']'
		}

		depth := 0
		in_string := false
		escaped := false

		for *index < len(raw) {

			current := raw[*index]
			*index++

			if current == '\\' && in_string == true && escaped == false {
				escaped = true
				continue
			}

			if current == '"' && escaped == false {
				in_string = !in_string
			}

			if in_string == false {

				if current == character {
					depth++
				} else if current == close_token {

					depth--

					if depth == 0 {
						return raw[start:*index], nil
					}

				}

			}

			escaped = false

		}

		return "", fmt.Errorf("commands.Parse: Unterminated nested value.")

	}

	for *index < len(raw) && unicode.IsSpace(rune(raw[*index])) == false {
		*index++
	}

	return raw[start:*index], nil

}

func skipWhitespace(raw string, index *int) {

	for *index < len(raw) && unicode.IsSpace(rune(raw[*index])) {
		*index++
	}

}
