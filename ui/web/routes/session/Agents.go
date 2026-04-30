package session

import "exocomp/ui/web/handlers"
import "exocomp/schemas"
import "exocomp/tools"
import "exocomp/types"
import "encoding/json"
import "net/http"
import "sort"
import "strconv"

func Agents(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		tool := session.GetTool("agents.List")

		if tool != nil {

			agent_tool, ok := tool.(*tools.Agents)

			if ok == true {

				agent_names := make([]string, 0)

				for name, _ := range agent_tool.Agents {
					agent_names = append(agent_names, name)
				}

				sort.Strings(agent_names)

				agents := make([]schemas.Agent, 0)

				agents = append(agents, schemas.Agent{
					Name:        session.Agent.Name,
					Type:        session.Agent.Type,
					Model:       session.Agent.Model,
					Temperature: session.Agent.Temperature,
					Messages:    session.Agent.Messages,
					Programs:    session.Agent.Programs,
					Tools:       session.Agent.Tools,
					Sandbox:     session.Agent.Sandbox,
				})

				for _, name := range agent_names {

					agent, ok := agent_tool.Agents[name]

					if ok == true {

						agents = append(agents, schemas.Agent{
							Name:        agent.Name,
							Type:        agent.Type,
							Model:       agent.Model,
							Temperature: agent.Temperature,
							Messages:    agent.Messages,
							Programs:    agent.Programs,
							Tools:       agent.Tools,
							Sandbox:     agent.Sandbox,
						})

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
