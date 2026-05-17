package session

import "exocomp/ui/web/handlers"
import "exocomp/tools"
import "exocomp/types"
import "encoding/json"
import "net/http"
import "strconv"

func Agents(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		tool := session.GetTool("agents.List")

		if tool != nil {

			agent_tool, ok := tool.(*tools.Agents)

			if ok == true {

				agent_names := agent_tool.GetNames()

				agents := make([]*types.Agent, 0)
				agents = append(agents, session.Agent)

				for _, name := range agent_names {

					agent := agent_tool.GetAgent(name)

					if agent != nil {
						agents = append(agents, agent)
					}

				}

				response_payload, err0 := json.MarshalIndent(agents, "", "\t")

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
			handlers.NotFound(session, request, response)
		}

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
