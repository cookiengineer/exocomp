package tty

import "exocomp/adapters"
import "exocomp/engine"
import "exocomp/tools"
import "exocomp/types"

func RestoreClient(session *engine.Session, agents []*types.Agent) *Client {

	if session.Agent.HasTools() {

		toolset := tools.Toolset(
			session.Config.Playground,
			session.Config.Sandbox,
			session.Config.Model,
			session.Config.URL,
			session.Config.Debug,
			session.Agent.AllowedPrograms,
			session.Agent.AllowedTools,
		)

		if len(toolset) > 0 {

			for _, tool := range toolset {
				session.SetTool(tool)
			}

		}

	}

	if session.Config.HasProvider(session.Agent.Model) {

		adapterset := adapters.Adapterset(
			session.Config.ResolveURL(session.Agent.Model, ""),
			session.Config.ResolveModel(session.Agent.Model),
		)

		if len(adapterset) > 0 {

			for _, adapter := range adapterset {
				session.SetAdapter(adapter)
			}

		}

	}

	tool := session.GetTool("agents.List")

	if tool != nil && len(agents) > 0 {

		agents_tool, ok := tool.(*tools.Agents)

		if ok == true {

			for _, agent := range agents {
				agents_tool.SetAgent(agent)
			}

		}

	}

	client := &Client{
		App:     NewApp(session),
		Session: session,
		Role:    "user",
	}

	client.bindEvents()

	return client

}
