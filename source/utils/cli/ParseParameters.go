package cli

import "encoding/json"
import "strconv"
import "unicode"

func seekNested(raw string, start int, start_token byte, close_token byte) []byte {

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

func seekString(raw string, start int, token byte) []byte {

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

func skipWhiteSpace(raw string, index *int) {

	for *index < len(raw) && unicode.IsSpace(rune(raw[*index])) {
		*index++
	}

}

func ParseParameters(raw string) map[string]any {

	parameters := make(map[string]any)
	index      := 0

	for index < len(raw) {

		skipWhiteSpace(raw, &index)

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

		skipWhiteSpace(raw, &index)

		value := make([]byte, 0)

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
			skipWhiteSpace(raw, &index)

		}

		if len(value) > 0 {
			parameters[key] = parseValue(value)
			index += len(value)
		}

	}

	return parameters

}

func parseValue(buffer []byte) any {

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

