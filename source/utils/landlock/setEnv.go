package landlock

import "strings"

func setEnv(environment []string, name string, value string) []string {

	prefix := name + "="
	result := make([]string, 0, len(environment)+1)

	for _, entry := range environment {

		if strings.HasPrefix(entry, prefix) == false {
			result = append(result, entry)
		}

	}

	result = append(result, prefix+value)

	return result

}
