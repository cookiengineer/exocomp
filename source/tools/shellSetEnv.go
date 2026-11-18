package tools

import "strings"

// shellSetEnv replaces or appends an environment variable.
func shellSetEnv(environment []string, name string, value string) []string {

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
