package routes

import "exocomp/types"
import "exocomp/ui/web/handlers"
import "net/http"

func AgentSettings(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodPost {

		// TODO: Modify session.Config.Name
		// TODO: Modify session.Config.Agent
		// TODO: Modify session.Config.Model
		// TODO: Modify session.Config.Temperature

		// TODO: Modify session.Agent instance with new settings

		handlers.NotFound(session, request, response)

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
