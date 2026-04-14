package routes

import "exocomp/ui/web/handlers"
import "exocomp/schemas"
import "exocomp/types"
import "encoding/json"
import "io"
import "net/http"
import "strconv"

func Chat(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodPost {

		content_type := request.Header.Get("Content-Type")

		if content_type == "application/json" {

			message := schemas.Message{}
			from    := len(session.Messages)

			request_payload, err0 := io.ReadAll(request.Body)

			if err0 == nil {

				err1 := json.Unmarshal(request_payload, message)

				if err1 == nil {

					err2 := session.SendChatRequest(message)

					if err2 == nil {

						messages := session.GetMessages(from)

						response_payload, err3 := json.MarshalIndent(messages, "", "\t")

						if err3 == nil {

							response.Header().Set("Content-Type", "application/json")
							response.Header().Set("Content-Length", strconv.Itoa(len(response_payload)))
							response.WriteHeader(http.StatusOK)
							response.Write(response_payload)

						} else {
							handlers.InternalServerError(session, request, response)
						}

					} else {
						handlers.InternalServerError(session, request, response)
					}

				} else {
					handlers.BadRequest(session, request, response)
				}

			} else {
				handlers.UnsupportedMediaType(session, request, response)
			}

		} else {
			handlers.UnsupportedMediaType(session, request, response)
		}

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
