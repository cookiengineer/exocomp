package sandbox

import "strings"

func unsetEnv(environment []string, name string) []string {

	prefix := name + "="
	result := make([]string, 0, len(environment))

	for _, entry := range environment {

		if strings.HasPrefix(entry, prefix) == false {
			result = append(result, entry)
		}

	}

	return result

}
