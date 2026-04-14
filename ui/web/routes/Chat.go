package routes

import "exocomp/ui/web/handlers"
import "exocomp/schemas"
import "exocomp/types"
import "encoding/json"
import "io"
import "net/http"

func Chat(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodPost {

		content_type := request.Header.Get("Content-Type")

		if content_type == "application/json" {

			message := schemas.Message{}
			amount  := len(session.Messages)

			bytes, err0 := io.ReadAll(request.Body)

			if err0 == nil {

				err1 := json.Unmarshal(bytes, message)

				if err1 == nil {

					err2 := session.SendChatRequest(message)

					if err2 == nil {

						// TODO: Response should be the new history

					} else {

						// TODO

					}

				} else {
					// TODO
				}

			} else {
				// TODO: Invalid Payload
			}

		} else {
			// TODO: Invalid payload
		}

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
