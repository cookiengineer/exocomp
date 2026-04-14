package routes

import "exocomp/types"
import "exocomp/ui/web/handlers"
import "net/http"

func Agents(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		// TODO

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
