package markdown

import "strings"

func isHorizontalRule(line string) bool {

	trimmed := strings.ReplaceAll(strings.TrimSpace(line), " ", "")

	if len(trimmed) < 3 {
		return false
	}

	first := trimmed[0]

	if first != '-' && first != '*' && first != '_' {
		return false
	}

	for index := 1; index < len(trimmed); index++ {
		if trimmed[index] != first {
			return false
		}
	}

	return true

}
