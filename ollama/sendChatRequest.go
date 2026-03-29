package ollama

import "exocomp/schemas"
import "bytes"
import "encoding/json"
import "io"
import "fmt"

func sendChatRequest(session *Session) error {

	request_payload, err0 := json.Marshal(schemas.ChatRequest{
		Model:       session.Agent.Model,
		Temperature: session.Agent.Temperature,
		Messages:    session.Messages,
		Stream:      false,
		Tools:       session.Tools,
	})

	if err0 == nil {

		endpoint := session.Config.ResolvePath("/api/chat")

		response, err1 := session.Client.Post(
			endpoint.String(),
			"application/json",
			bytes.NewReader(request_payload),
		)

		if err1 == nil && response.StatusCode == 200 {

			response_payload, err2 := io.ReadAll(response.Body)

			if err2 == nil {

				var response schemas.ChatResponse

				err3 := json.Unmarshal(response_payload, &response)

				if err3 == nil {
					return processChatResponse(session, response.Message)
				} else {
					return err3
				}

			} else {
				return err2
			}

		} else if err1 == nil && response.StatusCode == 404 {
			return fmt.Errorf("Ollama model %s not found", session.Config.Model)
		} else {
			return err1
		}

	} else {
		return err0
	}

}
