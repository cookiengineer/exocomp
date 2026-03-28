package ollama

import "fmt"
import "exocomp/parsers"
import "strings"

func processChatResponse(session *Session, response Message) error {

	if response.Role == "assistant" {

		session.mutex.Lock()
		session.Messages = append(session.Messages, &response)
		session.mutex.Unlock()

		// TODO: response.ToolCalls is set

		tools := parsers.ParseTools(
			session.Config.Agent,
			session.Config.Sandbox,
			session.Config.Tools,
			session.Config.Programs,
			response.Content,
		)

		if len(tools) > 0 {

			for _, tool := range tools {

				result, err0 := tool.Call()

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
