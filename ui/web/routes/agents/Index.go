package agents

import "exocomp/ui/web/handlers"
import "exocomp/tools"
import "exocomp/types"
import "encoding/json"
import "net/http"
import "strconv"

func Index(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		tool := session.GetTool("agents.List")

		if tool != nil {

			agent_tool, ok := tool.(*tools.Agents)

			if ok == true {

				agents := make(map[string]string)

				for name, agent := range agent_tool.Agents {
					agents[name] = agent.Type.String()
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
