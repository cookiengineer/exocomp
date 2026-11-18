package tty

import "exocomp/types"

func Run(agent *types.Agent, config *types.Config, role string) {

	client := NewClient(agent, config)

	if role != "" {
		client.SetRole(role)
	}

	client.Init()

}
