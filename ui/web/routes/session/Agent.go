package session

import "exocomp/ui/web/handlers"
import "exocomp/schemas"
import "exocomp/types"
import "encoding/json"
import "net/http"
import "strconv"

func Agent(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		response_payload, err0 := json.MarshalIndent(schemas.Agent{
			Name:            session.Agent.Name,
			Type:            session.Agent.Type,
			Model:           session.Agent.Model,
			Temperature:     session.Agent.Temperature,
			Messages:        session.Agent.Messages,
			AllowedPrograms: session.Agent.AllowedPrograms,
			AllowedTools:    session.Agent.AllowedTools,
			Sandbox:         session.Agent.Sandbox,
		}, "", "\t")

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
