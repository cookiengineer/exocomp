package tools

import "exocomp/types"
import "strings"

func getAgentWorkReport(agent *types.Agent) string {

	if agent != nil {

		for m := len(agent.Messages) - 1; m >= 0; m-- {

			message := agent.Messages[m]

			if message.Role == "tool" && strings.HasPrefix(message.Content, "agents.Quit: Work Report\n") {
				return strings.TrimSpace(strings.TrimPrefix(message.Content, "agents.Quit: Work Report\n"))
			}

		}

	}

	return ""

}
