package lint

import "strings"

func isCommentLine(line string) bool {

	trimmed := strings.TrimSpace(line)

	return strings.HasPrefix(trimmed, "//") == true || strings.HasPrefix(trimmed, "/*") == true

}
