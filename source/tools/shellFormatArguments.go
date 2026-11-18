package tools

import "fmt"
import "strings"

// shellFormatArguments formats an argument list for tool output.
func shellFormatArguments(arguments []string) string {

	if len(arguments) > 0 {

		quoted := make([]string, 0, len(arguments))

		for _, argument := range arguments {
			quoted = append(quoted, fmt.Sprintf("%q", argument))
		}

		return "[" + strings.Join(quoted, " ") + "]"

	} else {
		return "[]"
	}

}
