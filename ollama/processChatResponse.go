package ollama

import "exocomp/config"
import "strings"

func processChatResponse(session *Session, response Message) error {

	if response.Role == "assistant" {

		session.mutex.Lock()
		session.Messages = append(session.Messages, &response)
		session.mutex.Unlock()

		gadget := config.ParseGadget(response.Content)

		if gadget != nil && session.Config.IsAllowedGadget(gadget.Type.String()) == true {

			result, err0 := gadget.Call(session.Config)

			if err0 == nil {

				session.mutex.Lock()
				session.Messages = append(session.Messages, &Message{
					Role:    "tool",
					Content: strings.TrimSpace(result),
				})
				session.mutex.Unlock()

				return sendChatRequest(session)

			} else {
				return err0
			}

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
