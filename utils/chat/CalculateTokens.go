package chat

import "exocomp/schemas"

func CalculateTokens(messages []*schemas.Message) int {

	// estimation is typically between ~3.5 and 4.0 for most LLMs
	chars_per_token := float64(4.0)
	result := int(0)

	for _, message := range messages {

		if message == nil {
			continue
		}

		result += 4 // {"role","content","tool_name"}

		result += estimateStringTokens(message.Role, chars_per_token)

		if message.Content != "" {
			result += estimateStringTokens(message.Content, chars_per_token)
		}

		if message.ToolName != "" {
			result += estimateStringTokens(message.ToolName, chars_per_token)
		}

		if len(message.ToolCalls) > 0 {

			for _, tool_call := range message.ToolCalls {

				result += 8 // {"tool_call":{"id","type","name","arguments"}}

				result += estimateStringTokens(tool_call.ID,            chars_per_token)
				result += estimateStringTokens(tool_call.Type,          chars_per_token)
				result += estimateStringTokens(tool_call.Function.Name, chars_per_token)

				if len(tool_call.Function.ArgumentsRaw) > 0 {
					result += estimateBytesTokens(tool_call.Function.ArgumentsRaw, chars_per_token)
				}

			}

		}

	}

	return result

}

