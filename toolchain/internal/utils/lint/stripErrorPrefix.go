package lint

import "strings"

// stripErrorPrefix removes a leading component prefix.
func stripErrorPrefix(message string) string {

	index := strings.Index(message, ": ")

	if index < 0 {
		return message
	} else {

		head := message[:index]

		if prefixedErrorPattern.MatchString(head+":") == true {
			return message[index+2:]
		} else {
			return message
		}

	}

}
