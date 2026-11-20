package agents

import "sort"

func (tool *Agents) GetNames() []string {

	result := make([]string, 0)

	tool.mutex.Lock()
	for name, _ := range tool.contents {
		result = append(result, name)
	}
	tool.mutex.Unlock()

	sort.Strings(result)

	return result

}
