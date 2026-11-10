package jsonl

import "exocomp/engine"
import "exocomp/tools"

func RestoreClient(session *engine.Session) *Client {

	renderer := NewRenderer(session)

	session.InstallTools()
	session.InstallAdapters()

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
