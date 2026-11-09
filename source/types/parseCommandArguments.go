package types

import utils_strings "exocomp/utils/strings"

func parseCommandArguments(raw string) map[string]any {

	arguments := make(map[string]any)
	index := 0

	for index < len(raw) {

		skipWhitespace(raw, &index)

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

		skipWhitespace(raw, &index)

		value := make([]byte, 0)

		switch raw[index] {
		case '{':
			value = utils_strings.SeekNested(raw, index, '{', '}')

		case '[':
			value = utils_strings.SeekNested(raw, index, '[', ']')

		case '\'':
			value = utils_strings.Seek(raw, index, '\'')

		case '"':
			value = utils_strings.Seek(raw, index, '"')

		default:
			skipWhitespace(raw, &index)

		}

		if len(value) > 0 {
			arguments[key] = parseCommandArgument(value)
			index += len(value)
		}

	}

	return arguments

}

