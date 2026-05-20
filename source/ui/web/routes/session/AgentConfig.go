package session

import "exocomp/ui/web/handlers"
import "exocomp/tools"
import "exocomp/types"
import "encoding/json"
import "net/http"
import "strconv"

func AgentConfig(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		name := request.PathValue("name")

		if name == session.Config.Name {

			if session.Config != nil {

				response_payload, err0 := json.MarshalIndent(session.Config, "", "\t")

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

		} else if name != "" {

			tool := session.GetTool("agents.List")

			if tool != nil {

				agent_tool, ok0 := tool.(*tools.Agents)

				if ok0 == true {

					agent := agent_tool.GetAgent(name)

					if agent != nil {

						response_payload, err1 := json.MarshalIndent(&types.Config{
							Name:        agent.Name,
							Role:        agent.Role,
							Model:       agent.Model,
							Prompt:      agent.Prompt,
							Temperature: agent.Temperature,
							Playground:  session.Config.Playground,
							Sandbox:     agent.Sandbox,
							URL:         session.Config.URL,
							Debug:       session.Config.Debug,
						}, "", "\t")

						if err1 == nil {

							response.Header().Set("Content-Type", "application/json")
							response.Header().Set("Content-Length", strconv.Itoa(len(response_payload)))
							response.WriteHeader(http.StatusOK)
							response.Write(response_payload)

						} else {
							handlers.InternalServerError(session, err1, request, response)
						}

					} else {
						handlers.NotFound(session, request, response)
					}

				} else {
					handlers.NotFound(session, request, response)
				}

			} else {
				handlers.NotFound(session, request, response)
			}

		} else {
			handlers.BadRequest(session, nil, request, response)
		}

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
