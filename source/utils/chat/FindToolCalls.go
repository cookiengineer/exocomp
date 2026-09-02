package chat

import "exocomp/schemas"
import "fmt"
import "strings"

func FindToolCalls(message *schemas.Message, search []string) ([]schemas.ToolCall, error) {

	result := make([]schemas.ToolCall, 0)

	for _, tool_call := range message.ToolCalls {

		toolname, err1 := tool_call.GetName()

		if err1 == nil {

			matches := false

			for _, keyword := range search {

				if toolname == keyword {
					matches = true
					break
				}

			}

			if matches == true {
				result = append(result, tool_call)
			}

		}

	}

	if len(result) > 0 {
		return result, nil
	} else {
		return result, fmt.Errorf("No tool calls matching \"%s\" found.", strings.Join(search, "\",\""))
	}

}
