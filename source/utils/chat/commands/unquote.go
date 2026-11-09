package commands

import "strconv"

func unquote(value string) string {

	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {

		unquoted, err := strconv.Unquote(value)

		if err == nil {
			return unquoted
		}

	}

	if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
		return value[1 : len(value)-1]
	}

	return value

}
