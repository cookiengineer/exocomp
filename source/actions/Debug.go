package actions

import "exocomp/engine"
import "exocomp/types"
import "exocomp/ui/tty"

func Debug(session *engine.Session, agents []*types.Agent, role string) {

	client := tty.RestoreClient(session, agents)

	if role != "" {
		client.SetRole(role)
	}

	client.Init()

}
