package yaml

import "strings"

func isSequenceItem(text string) bool {

	if strings.HasPrefix(text, "- ") == true {
		return true
	} else {
		return text == "-"
	}
}
