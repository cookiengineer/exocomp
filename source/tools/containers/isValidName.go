package containers

import "strings"

func (tool *Containers) isValidName(name string) bool {

	if name == "" || len(name) > 64 {
		return false
	}

	for index, character := range name {

		if character >= 'a' && character <= 'z' {
			continue
		} else if character >= 'A' && character <= 'Z' {
			continue
		} else if character >= '0' && character <= '9' {
			continue
		} else if index > 0 && strings.ContainsRune("._-", character) == true {
			continue
		} else {
			return false
		}

	}

	return true

}
