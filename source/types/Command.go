package types

import "encoding/json"
import "fmt"
import "strconv"
import "strings"
import "unicode"

func parse_arguments(raw string) map[string]any {

	arguments := make(map[string]any)
	index     := 0

	for index < len(raw) {

		skip_whitespace(raw, &index)

		if index >= len(raw) {
			break
		}

		// key
		key_start := index

		for index < len(raw) && raw[index] != '=' {
			index++
		}

		if index >= len(raw) {
			break
		}

		key := raw[key_start:index]

		// Skip '='
		index++

		skip_whitespace(raw, &index)

		value := make([]byte, 0)

		switch raw[index] {
		case '{':
			value = seek_nested(raw, index, '{', '}')

		case '[':
			value = seek_nested(raw, index, '[', ']')

		case '\'':
			value = seek_string(raw, index, '\'')

		case '"':
			value = seek_string(raw, index, '"')

		default:
			skip_whitespace(raw, &index)

		}

		if len(value) > 0 {
			arguments[key] = parse_value(value)
			index += len(value)
		}

	}

	return arguments

}

func parse_value(buffer []byte) any {

	var value any

	err := json.Unmarshal(buffer, &value)

	if err == nil {

		return value

	} else {

		if len(buffer) >= 2 && (buffer[0] == '-' || buffer[0] == '+') {

			num, err1 := strconv.ParseInt(string(buffer), 10, 64);

			if err1 == nil {
				return num
			}

		}

		if len(buffer) >= 2 && buffer[0] == '"' && buffer[len(buffer)-1] == '"' {
			return string(buffer[1:len(buffer)-1])
		}

		if len(buffer) >= 2 && buffer[0] == '\'' && buffer[len(buffer)-1] == '\'' {
			return string(buffer[1:len(buffer)-1])
		}

		if string(buffer) == "true" {
			return true
		}

		if string(buffer) == "false" {
			return false
		}

		if string(buffer) == "null" {
			return nil
		}

		num1, err1 := strconv.ParseUint(string(buffer), 10, 64)

		if err1 == nil {
			return num1
		}

		num2, err2 := strconv.ParseFloat(string(buffer), 64)

		if err2 == nil {
			return num2
		}

	}

	return nil

}

func seek_nested(raw string, start int, start_token byte, close_token byte) []byte {

	result    := []byte{start_token}
	depth     := 0
	in_string := false
	escaped   := false

	for index := start + 1; index < len(raw); index++ {

		result = append(result, byte(raw[index]))

		chr := raw[index]

		if chr == '\\' && in_string && escaped == false {
			escaped = true
			continue
		}

		if chr == '"' && escaped == false {
			in_string = !in_string
		}

		if in_string == false {

			if chr == start_token {

				depth++

			} else if chr == close_token {

				depth--

				if depth == 0 {
					break
				}

			}

		}

		escaped = false

	}

	return result

}

func seek_string(raw string, start int, token byte) []byte {

	result  := []byte{token}
	escaped := false

	for index := start + 1; index < len(raw); index++ {

		result = append(result, byte(raw[index]))

		chr := raw[index]

		if chr == '\\' && escaped == false {
			escaped = true
			continue
		}

		if chr == token && escaped == false {
			break
		}

		escaped = false

	}

	return result

}

func skip_whitespace(raw string, index *int) {

	for *index < len(raw) && unicode.IsSpace(rune(raw[*index])) {
		*index++
	}

}


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

			command.Name      = strings.TrimSpace(name)
			command.Method    = strings.TrimSpace(command.Name[strings.LastIndex(command.Name, ".")+1:])
			command.Arguments = parse_arguments(strings.TrimSpace(prompt[1+len(name)+1:]))

			return nil

		} else {
			return fmt.Errorf("invalid command tool name \"%s\"", name)
		}

	} else {
		return fmt.Errorf("invalid command prompt syntax")
	}

}
