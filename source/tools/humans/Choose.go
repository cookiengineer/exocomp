package humans

import "exocomp/types"
import "context"
import "strings"

func (tool *Humans) Choose(text string, options []string, multiple bool) (string, error) {

	question := &types.Question{
		Type:     "Choose",
		Question: strings.TrimSpace(text),
		Options:  options,
		Multiple: multiple,
		Answer:   "",
	}

	ctx, cancel := context.WithTimeout(context.Background(), tool.Timeout)
	state := &questionState{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}

	tool.mutex.Lock()
	tool.contents[question.Question] = question
	tool.states[question.Question] = state
	tool.mutex.Unlock()

	go watchQuestion(tool, question.Question, state)

	return tool.Await(question.Question)

}
