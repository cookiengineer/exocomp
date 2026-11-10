package tty

import "exocomp/engine"
import "exocomp/types"

func RestoreClient(session *engine.Session, agents []*types.Agent) *Client {

	session.InstallTools()
	session.InstallAdapters()
	session.InstallAgents(agents)

	client := &Client{
		App:     NewApp(session),
		Session: session,
		Role:    "user",
	}

	client.bindEvents()

	return client

}
