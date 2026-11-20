package agents

import "exocomp/types"
import "sort"

func (tool *Agents) GetAgents() []*types.Agent {

	tool.mutex.Lock()

	result := make([]*types.Agent, 0, len(tool.contents))

	for _, agent := range tool.contents {
		result = append(result, agent.Clone())
	}

	tool.mutex.Unlock()

	sort.Slice(result, func(a int, b int) bool {
		return result[a].Name < result[b].Name
	})

	return result

}
