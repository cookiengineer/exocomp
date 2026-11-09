package commands

import utils_strings "exocomp/utils/strings"
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

func ParseRawArguments(raw string) (map[string]string, error) {

	result := make(map[string]string)
	index := 0

	for index < len(raw) {

		utils_strings.SkipWhitespace(raw, &index)

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

		utils_strings.SkipWhitespace(raw, &index)

		value, err := scanValue(raw, &index)

		if err == nil {
			result[key] = value
		} else {
			return nil, err
		}

	}

	return result, nil

}
