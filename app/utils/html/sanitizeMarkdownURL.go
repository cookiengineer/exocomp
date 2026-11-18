package html

import "strings"

func sanitizeMarkdownURL(value string) string {

	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return ""
	}

	lower := strings.ToLower(trimmed)

	if strings.HasPrefix(lower, "javascript:") == true || strings.HasPrefix(lower, "data:") == true || strings.HasPrefix(lower, "vbscript:") == true {
		return ""
	}

	return trimmed

}
