package humans

import "exocomp/types"
import "sort"

func (tool *Humans) Snapshot() []*types.Question {

	tool.mutex.RLock()

	result := make([]*types.Question, 0)

	for _, question := range tool.contents {

		if question.Answer == "" {
			result = append(result, question)
		}

	}

	tool.mutex.RUnlock()

	sort.Slice(result, func(a int, b int) bool {
		return result[a].Question < result[b].Question
	})

	return result

}
