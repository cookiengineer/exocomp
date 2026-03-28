package ollama

import "exocomp/schemas"
import "fmt"
import "strings"

func processChatResponse(session *Session, response schemas.Message) error {

	if response.Role == "assistant" {

		session.mutex.Lock()
		session.Messages = append(session.Messages, response)
		session.mutex.Unlock()

		if len(response.ToolCalls) > 0 {

			for _, tool_call := range response.ToolCalls {

				name,      err0 := tool_call.Function.Tool()
				method,    err1 := tool_call.Function.Method()
				arguments, err2 := tool_call.Function.Arguments()

				if err0 == nil && err1 == nil && err2 == nil {

					tool := session.GetTool(name)

					if tool != nil {

						result, err0 := tool.Call(method, arguments)

						if err0 == nil {

							session.mutex.Lock()
							session.Messages = append(session.Messages, schemas.Message{
								Role:      "tool",
								Content:   strings.TrimSpace(result),
								ToolName:  name + "." + method,
								ToolCalls: []schemas.ToolCall{tool_call},
							})
							session.mutex.Unlock()

						} else {

							session.mutex.Lock()
							session.Messages = append(session.Messages, schemas.Message{
								Role:      "tool",
								Content:   fmt.Sprintf("Error: %s", strings.TrimSpace(err0.Error())),
								ToolName:  name + "." + method,
								ToolCalls: []schemas.ToolCall{tool_call},
							})
							session.mutex.Unlock()

						}

					}

				}

			}

			return sendChatRequest(session)

		} else {
			return nil
		}

	} else {

		session.mutex.Lock()
		session.Messages = append(session.Messages, response)
		session.mutex.Unlock()

		return nil

	}

}
