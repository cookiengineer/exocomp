package components

import "exocomp/schemas"
import "encoding/json"
import "strings"

func toolCallContent(tool_call schemas.ToolCall) string {

	arguments, err := tool_call.GetArguments()

	if err == nil {

		bytes, err0 := json.MarshalIndent(arguments, "", "\t")

		if err0 == nil {
			return "```json\n" + string(bytes) + "\n```"
		}

	}

	raw := strings.TrimSpace(string(tool_call.Function.ArgumentsRaw))

	if raw == "" {
		return ""
	} else {
		return "```json\n" + raw + "\n```"
	}

}
