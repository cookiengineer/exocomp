package html

import "strings"

func needsMarkdownSpace(previous string, next string) bool {

	if previous == "" || next == "" {
		return false
	}

	last := previous[len(previous)-1]
	first := next[0]

	if last == ' ' || last == '\n' || last == '\t' {
		return false
	}

	if first == ' ' || first == '\n' || first == '\t' {
		return false
	}

	return strings.ContainsRune(".,;:!?)]}", rune(first)) == false

}
