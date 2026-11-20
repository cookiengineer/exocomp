package agents

import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"

func (tool *Agents) SetAgent(agent *types.Agent) bool {

	sanitized_name := utils_fmt.FormatAgentName(agent.Name)

	if sanitized_name == agent.Name {

		tool.mutex.Lock()
		_, ok := tool.contents[sanitized_name]
		tool.mutex.Unlock()

		if ok == false {

			tool.mutex.Lock()
			tool.contents[sanitized_name] = agent
			tool.mutex.Unlock()

			return true

		}

	}

	return false

}
