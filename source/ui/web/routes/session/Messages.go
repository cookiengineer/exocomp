package session

import "exocomp/ui/web/handlers"
import "exocomp/types"
import "encoding/json"
import "net/http"
import "strconv"

func Messages(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		messages               := session.GetMessages(0)
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
		handlers.MethodNotAllowed(session, request, response)
	}

}
