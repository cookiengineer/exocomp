package markdown

import "strings"

func isHorizontalRule(line string) bool {

	trimmed := strings.ReplaceAll(strings.TrimSpace(line), " ", "")

	if len(trimmed) >= 3 {

		first := trimmed[0]

		if first == '-' || first == '*' || first == '_' {

			for index := 1; index < len(trimmed); index++ {
				if trimmed[index] != first {
					return false
				}
			}

			return true

		} else {
			return false
		}

	} else {
		return false
	}

}
