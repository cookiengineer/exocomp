package markdown

import "strings"

func getTokenPosition(value string, remaining string) int {
	return len(strings.TrimSpace(value)) - len(remaining)
}
