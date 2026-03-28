package ollama

import "exocomp/schemas"
import "exocomp/tools"
import "fmt"
import "strings"

func processChatResponse(session *Session, response schemas.Message) error {

	if response.Role == "assistant" {

		session.mutex.Lock()
		session.Messages = append(session.Messages, &response)
		session.mutex.Unlock()

		if len(response.ToolCalls) > 0 {

			for _, tool_call := range response.ToolCalls {

				function  := strings.TrimSpace(tool_call.Function)
				arguments := tool_call.Function.Arguments()

				var tool *tools.Tool = nil

				// TODO: Implement bugs.* tool
				// TODO: Implement features.* tool

				if strings.HasPrefix(function, "files.") {
					tool = session.RequestTool("files")
				} else if strings.HasPrefix(function, "notes.") {
					tool = session.RequestTool("notes")
				} else if strings.HasPrefix(function, "programs.") {
					tool = session.RequestTool("programs")
				}

				if tool != nil {

					result, err0 := tool.Call(method, arguments)

					if err0 == nil {

						session.mutex.Lock()
						session.Messages = append(session.Messages, &Message{
							Role:    "tool",
							Content: strings.TrimSpace(result),
						})
						session.mutex.Unlock()

					} else {

						session.mutex.Lock()
						session.Messages = append(session.Messages, &Message{
							Role:    "tool",
							Content: fmt.Sprintf("Error: %s", strings.TrimSpace(err0.Error())),
						})
						session.mutex.Unlock()

					}

				}

			}

			return sendChatRequest(session)

		} else {
			return nil
		}

	} else {

		session.mutex.Lock()
		session.Messages = append(session.Messages, &response)
		session.mutex.Unlock()

		return nil

	}

}
