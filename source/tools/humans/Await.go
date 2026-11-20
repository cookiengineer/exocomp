package humans

import "fmt"

// NOTE: Await blocks until the human answered the question. This mirrors
// agents.Await deliberately: the planner model would otherwise poll in a hot
// loop and blow up its limited context window with identical "still waiting"
// tool messages. The lifecycle guarantees the done channel always closes
// (answer received -> Answer() -> close, or timeout -> watchQuestion ->
// close), so Await never blocks forever.
func (tool *Humans) Await(reference string) (string, error) {

	tool.mutex.Lock()
	state, waiting := tool.states[reference]
	tool.mutex.Unlock()

	if waiting == true {
		<-state.done
	}

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	question, ok := tool.contents[reference]

	if ok == true {

		if question.Answer != "" {
			return fmt.Sprintf("humans.%s: Answer for Question \"%s\":\n===\n%s\n===", question.Type, question.Question, question.Answer), nil
		} else {
			return "", fmt.Errorf("humans.%s: Question \"%s\" was never answered!", question.Type, question.Question)
		}

	} else {
		return "", fmt.Errorf("humans.Await: Question \"%s\" was never asked!", reference)
	}

}
