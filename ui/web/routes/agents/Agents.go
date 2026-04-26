package agents

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

				agents := make(map[string]schemas.Agent)

				for _, name := range agent_names {

					agent,    ok1 := agent_tool.Agents[name]
					messages, ok2 := agent_tool.Chats[name]

					if ok1 == true && ok2 == true {

						agents[name] = schemas.Agent{
							Name:        agent.Name,
							Type:        agent.Type,
							Model:       agent.Model,
							Programs:    agent.Programs,
							Temperature: agent.Temperature,
							Tools:       agent.Tools,
							Messages:    messages,
						}

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
