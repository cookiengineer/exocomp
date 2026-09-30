package cli

import "encoding/json"
import "strconv"
import "strings"

func ParseParameters(raw string) map[string]any {

	parameters := make(map[string]any)
	index := 0

	for index < len(raw) {

		index = skipWhitespace(raw, index)

		if index >= len(raw) {
			break
		}

		key_start := index

		for index < len(raw) && raw[index] != '=' {
			index++
		}

		if index >= len(raw) {
			break
		}

		key := raw[key_start:index]

		index++

		index = skipWhitespace(raw, index)

		value := ""

		if index < len(raw) {

			switch raw[index] {

			case '{':
				value = seekNested(raw, index, '{', '}')
			case '[':
				value = seekNested(raw, index, '[', ']')
			case '\'':
				value = seekString(raw, index, '\'')
			case '"':
				value = seekString(raw, index, '"')
			default:
				index = skipWhitespace(raw, index)
				value = seekToken(raw, index)

			}

		}

		if len(value) > 0 {
			parameters[key] = parseValue(value)
			index += len(value)
		}

	}

	return parameters

}

func seekNested(raw string, start int, start_token byte, close_token byte) string {

	result := string(start_token)
	depth := 1
	in_string := false
	escaped := false

	for index := start + 1; index < len(raw); index++ {

		chr := raw[index]
		result += string(chr)

		if chr == '\\' && in_string == true && escaped == false {
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

func seekString(raw string, start int, token byte) string {

	result := string(token)
	escaped := false

	for index := start + 1; index < len(raw); index++ {

		chr := raw[index]
		result += string(chr)

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

func seekToken(raw string, start int) string {

	index := start

	for index < len(raw) && raw[index] != ' ' && raw[index] != '\t' && raw[index] != '\n' && raw[index] != '\r' {
		index++
	}

	return raw[start:index]

}

func skipWhitespace(raw string, index int) int {

	for index < len(raw) && isWhitespace(raw[index]) == true {
		index++
	}

	return index

}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

func parseValue(buffer string) any {

	var value any

	if err := json.Unmarshal([]byte(buffer), &value); err == nil {
		return value
	}

	if len(buffer) >= 2 && (buffer[0] == '-' || buffer[0] == '+') {
		if num, err := strconv.Atoi(buffer); err == nil {
			return num
		}
	}

	if len(buffer) >= 2 && buffer[0] == '"' && buffer[len(buffer)-1] == '"' {
		return buffer[1 : len(buffer)-1]
	}

	if len(buffer) >= 2 && buffer[0] == '\'' && buffer[len(buffer)-1] == '\'' {
		return buffer[1 : len(buffer)-1]
	}

	if buffer == "true" {
		return true
	}

	if buffer == "false" {
		return false
	}

	if buffer == "null" {
		return nil
	}

	trimmed := strings.TrimSpace(buffer)

	if num, err := strconv.Atoi(trimmed); err == nil {
		return num
	}

	if num, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return num
	}

	return nil

}
