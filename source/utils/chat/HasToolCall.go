package chat

import "exocomp/schemas"

func HasToolCall(message *schemas.Message, toolname string) bool {

	if message.Role == "assistant" {

		if len(message.ToolCalls) > 0 {

			found := false

			for _, toolcall := range message.ToolCalls {

				if toolcall.Function.Name == toolname {
					found = true
					break
				}

			}

			return found

		}

	}

	return false

}
