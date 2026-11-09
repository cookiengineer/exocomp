package types

import utils_strings "exocomp/utils/strings"

func parseCommandArguments(raw string) map[string]any {

	arguments := make(map[string]any)
	index := 0

	for index < len(raw) {

		utils_strings.SkipWhitespace(raw, &index)

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

		utils_strings.SkipWhitespace(raw, &index)

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
			utils_strings.SkipWhitespace(raw, &index)

		}

		if len(value) > 0 {

			arguments[key] = parseCommandArgument(value)
			index += len(value)

		}

	}

	return arguments

}
