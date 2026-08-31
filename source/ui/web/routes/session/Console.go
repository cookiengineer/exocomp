package session

import "exocomp/engine"
import "exocomp/ui/web/handlers"
import "encoding/json"
import "net/http"
import "strconv"

func Console(session *engine.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		if session.Console != nil {

			messages               := session.Console.GetMessages(0)
			response_payload, err0 := json.MarshalIndent(messages, "", "\t")

			if err0 == nil {

				response.Header().Set("Content-Type", "application/json")
				response.Header().Set("Content-Length", strconv.Itoa(len(response_payload)))
				response.WriteHeader(http.StatusOK)
				response.Write(response_payload)

			} else {
				handlers.InternalServerError(session, err0, request, response)
			}

		} else {
			handlers.NotFound(session, request, response)
		}

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
