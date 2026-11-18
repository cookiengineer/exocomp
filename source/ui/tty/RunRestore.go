package tty

import "exocomp/engine"
import "exocomp/types"

func RunRestore(session *engine.Session, agents []*types.Agent, role string) {

	client := RestoreClient(session, agents)

	if role != "" {
		client.SetRole(role)
	}

	client.Init()

}
