package humans

import "fmt"

func (tool *Humans) Answer(question string, answer string) (string, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	answered := false

	for _, other := range tool.contents {

		if other.Question == question {

			other.Answer = answer

			state, is_waiting := tool.states[question]

			if is_waiting == true {

				delete(tool.states, question)
				state.cancel()
				close(state.done)

			}

			answered = true
			break

		}

	}

	if answered == true {
		return fmt.Sprintf("humans.Answer: Answer for Question \"%s\" is:\n===\n%s\n===", question, answer), nil
	} else {
		return "", fmt.Errorf("humans.Answer: Question \"%s\" was never asked!", question)
	}

}
