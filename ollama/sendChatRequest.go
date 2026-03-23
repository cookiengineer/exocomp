package ollama

import "bytes"
import "encoding/json"
import "io"

func sendChatRequest(session *Session) error {

	request_payload, err0 := json.Marshal(ChatRequest{
		Model:    session.Config.Model,
		Messages: session.Messages,
		Stream:   false,
	})

	if err0 == nil {

		endpoint := session.Config.ResolvePath("/api/chat")

		response, err1 := session.Client.Post(
			endpoint.String(),
			"application/json",
			bytes.NewReader(request_payload),
		)

		if err1 == nil {

			response_payload, err2 := io.ReadAll(response.Body)

			if err2 == nil {

				var response ChatResponse

				err3 := json.Unmarshal(response_payload, &response)

				if err3 == nil {
					return processChatResponse(session, *response.Message)
				} else {
					return err3
				}

			} else {
				return err2
			}

		} else {
			return err1
		}

	} else {
		return err0
	}

}
