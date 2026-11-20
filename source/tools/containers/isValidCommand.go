package containers

import "strings"

func (tool *Containers) isValidCommand(command string) bool {

	if command == "" {
		return false
	} else if strings.HasPrefix(command, "-") == true {
		return false
	} else if strings.ContainsAny(command, "\x00\n\r") == true {
		return false
	} else {
		return true
	}

}
