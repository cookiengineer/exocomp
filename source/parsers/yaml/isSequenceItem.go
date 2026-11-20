package yaml

import "strings"

func isSequenceItem(text string) bool {

	if strings.HasPrefix(text, "- ") {
		return true
	} else {
		return text == "-"
	}
}
