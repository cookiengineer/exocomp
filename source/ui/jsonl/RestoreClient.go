package jsonl

import "exocomp/adapters"
import "exocomp/engine"
import "exocomp/tools"

func RestoreClient(session *engine.Session) *Client {

	renderer := NewRenderer(session)

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

	client := &Client{
		Renderer: renderer,
		Session:  session,
		role:     "user",
	}

	tool := session.GetTool("agents.Quit")

	if tool != nil {

		agents_tool, ok := tool.(*tools.Agents)

		if ok == true {
			agents_tool.OnQuit = client.QuitHook
		}

	}

	return client

}
