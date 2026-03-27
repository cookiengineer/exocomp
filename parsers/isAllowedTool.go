package parsers

import "slices"

func isAllowedTool(tools []string, tool string) bool {
	return slices.Contains(tools, tool)
}
