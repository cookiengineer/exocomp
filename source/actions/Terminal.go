package actions

import "exocomp/types"
import "exocomp/ui/tty"

func Terminal(agent *types.Agent, config *types.Config, role string) {

	client := tty.NewClient(agent, config)

	if role != "" {
		client.SetRole(role)
	}

	client.Init()

}
