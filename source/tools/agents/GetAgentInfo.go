package agents

import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"

func (tool *Agents) GetAgentInfo(name string) *types.Agent {

	sanitized_name := utils_fmt.FormatAgentName(name)

	tool.mutex.Lock()
	agent, ok := tool.contents[sanitized_name]
	tool.mutex.Unlock()

	if ok == true {
		return agent.CloneInfo()
	} else {
		return nil
	}

}
