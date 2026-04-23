package routes

import "exocomp/ui/web/handlers"
import "exocomp/types"
import "encoding/json"
import "net/http"
import "strconv"

func Config(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		if session.Config != nil {

			response_payload, err0 := json.MarshalIndent(session.Config, "", "\t")

			if err0 == nil {

				response.Header().Set("Content-Type", "application/json")
				response.Header().Set("Content-Length", strconv.Itoa(len(response_payload)))
				response.WriteHeader(http.StatusOK)
				response.Write(response_payload)

			} else {
				handlers.InternalServerError(session, request, response)
			}

		} else {
			handlers.NotFound(session, request, response)
		}

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
