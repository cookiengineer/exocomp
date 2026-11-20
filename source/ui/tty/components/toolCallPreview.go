package components

import "exocomp/schemas"
import "fmt"
import "sort"
import "strings"

func toolCallPreview(tool_call schemas.ToolCall) string {

	raw := strings.TrimSpace(string(tool_call.Function.ArgumentsRaw))

	if raw == "" {
		return ""
	} else {

		arguments, err := tool_call.GetArguments()

		if err == nil {

			keys := make([]string, 0, len(arguments))

			for key := range arguments {
				keys = append(keys, key)
			}

			if len(keys) > 0 {

				sort.Strings(keys)
				key := keys[0]
				value := fmt.Sprintf("%v", arguments[key])

				return key + "=" + strings.Join(strings.Fields(value), " ")

			}

		}

		return strings.Join(strings.Fields(raw), " ")

	}

}
