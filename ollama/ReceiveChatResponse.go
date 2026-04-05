package ollama

import "encoding/json"
import "exocomp/schemas"
import "fmt"
import "strings"

func ReceiveChatResponse(session *Session, response schemas.Message) error {

	if response.Role == "assistant" {

		session.mutex.Lock()
		msg := &response
		session.Messages = append(session.Messages, msg)
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
							message := &schemas.Message{
								Role:     "tool",
								Content:  strings.TrimSpace(result),
								ToolName: name + "." + method,
							}
							session.Messages = append(session.Messages, message)
							session.mutex.Unlock()

						} else {

							session.mutex.Lock()
							message := &schemas.Message{
								Role:     "tool",
								Content:  fmt.Sprintf("Error: %s", strings.TrimSpace(err0.Error())),
								ToolName: name + "." + method,
							}
							session.Messages = append(session.Messages, message)
							session.mutex.Unlock()

						}

					} else {

						json_blob, _ := json.MarshalIndent(tool_call, "", "\t")

						session.mutex.Lock()
						message := &schemas.Message{
							Role:     "tool",
							Content:  strings.Join([]string{
								fmt.Sprintf("Error: %s", "Invalid Tool Call"),
								"",
								string(json_blob),
							}, "\n"),
							ToolName: name + "." + method,
						}
						session.Messages = append(session.Messages, message)
						session.mutex.Unlock()

					}

				}

			}

			return SendChatRequest(session)

		} else {
			return nil
		}

	} else {

		session.mutex.Lock()
		msg := &response
		session.Messages = append(session.Messages, msg)
		session.mutex.Unlock()

		return nil

	}

}
