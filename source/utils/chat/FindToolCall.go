package chat

import "exocomp/schemas"
import "fmt"

func FindToolCall(message *schemas.Message, search_toolname string, search_arguments map[string]any) (*schemas.ToolCall, error) {

	for _, tool_call := range message.ToolCalls {

		toolname,  err1 := tool_call.ToolName()
		arguments, err2 := tool_call.ToolArguments()

		if err1 == nil && err2 == nil {

			if toolname == search_toolname {

				if is_same_arguments(arguments, search_arguments) == true {
					return &tool_call, nil
				}

			}

		}

	}

	return nil, fmt.Errorf("No tool call for \"%s\" found.", search_toolname)

}
