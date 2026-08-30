package tools

import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"
import "context"
import "encoding/json"
import "fmt"
import "sync"
import "time"

var question_unique_id int64 = 0

type question_state struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

func watch_question(tool *Humans, id string, state *question_state) {

	<-state.ctx.Done()

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	_, still_waiting := tool.states[id]

	if still_waiting == true {
		delete(tool.states, id)
		close(state.done)
	}

}

type Humans struct {
	Timeout  time.Duration
	contents map[string]*types.Question
	states   map[string]*question_state
	mutex    *sync.Mutex
}

func NewHumans() *Humans {

	humans := &Humans{
		Timeout:  15 * time.Minute,
		contents: make(map[string]*types.Question),
		states:   make(map[string]*question_state),
		mutex:    &sync.Mutex{},
	}

	return humans

}

func (tool *Humans) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "Ask" {

		question, ok1 := arguments["question"].(string)

		if ok1 == true {
			return tool.Ask(utils_fmt.FormatMultiLine(question))
		} else {
			return "", fmt.Errorf("humans.%s: %s", method, "Invalid parameter \"question\" is not a string.")
		}

	} else if method == "Choose" {

		question, ok1 := arguments["question"].(string)
		raw,      ok2 := arguments["options"]
		multiple, _   := arguments["multiple"].(bool)

		if ok1 == true && ok2 == true {

			options := make([]string, 0)

			raw_options, ok4 := raw.([]interface{})

			if ok4 == true {

				for o, value := range raw_options {

					option, ok := value.(string)

					if ok == true {
						options = append(options, option)
					} else {
						return "", fmt.Errorf("humans.%s: %s", method, fmt.Sprintf("Invalid parameter \"options[%d]\" is not a string.", o))
					}

				}

			} else {
				return "", fmt.Errorf("humans.%s: %s", method, "Invalid parameter \"options\" is not an array of strings.")
			}

			return tool.Choice(utils_fmt.FormatMultiLine(question), options, multiple)

		} else if ok1 == true && ok2 == false {
			return "", fmt.Errorf("humans.%s: %s", method, "Invalid parameter \"options\" is not an array of strings.")
		} else if ok1 == false && ok2 == true {
			return "", fmt.Errorf("humans.%s: %s", method, "Invalid parameter \"question\" is not a string.")
		} else {
			return "", fmt.Errorf("humans.%s: Invalid parameters.", method)
		}

	} else {
		return "", fmt.Errorf("humans.%s: Invalid method.", method)
	}

}

func (tool *Humans) Answer(id string, answer string) error {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	question, ok := tool.contents[id]

	if ok == false {
		return fmt.Errorf("humans.Answer: Invalid question id.")
	}

	question.Answer = answer

	state, waiting := tool.states[id]

	if waiting == true {

		delete(tool.states, id)
		state.cancel()
		close(state.done)

	}

	return nil

}

func (tool *Humans) Ask(text string) (string, error) {

	question_unique_id++

	question := &types.Question{
		ID:       fmt.Sprintf("question-%d", question_unique_id),
		Type:     "Ask",
		Question: text,
		Options:  []string{},
		Multiple: false,
		Answer:   "",
	}

	ctx, cancel := context.WithTimeout(context.Background(), tool.Timeout)
	state := &question_state{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}

	tool.mutex.Lock()
	tool.contents[question.ID] = question
	tool.states[question.ID]  = state
	tool.mutex.Unlock()

	go watch_question(tool, question.ID, state)

	return tool.Await(question.ID)

}

func (tool *Humans) Choice(text string, options []string, multiple bool) (string, error) {

	question_unique_id++

	question := &types.Question{
		ID:       fmt.Sprintf("question-%d", question_unique_id),
		Type:     "Choice",
		Question: text,
		Options:  options,
		Multiple: multiple,
		Answer:   "",
	}

	ctx, cancel := context.WithTimeout(context.Background(), tool.Timeout)
	state := &question_state{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}

	tool.mutex.Lock()
	tool.contents[question.ID] = question
	tool.states[question.ID]  = state
	tool.mutex.Unlock()

	go watch_question(tool, question.ID, state)

	return tool.Await(question.ID)

}

func (tool *Humans) GetContent(id string) (any, error) {

	tool.mutex.Lock()
	content, ok := tool.contents[id]
	tool.mutex.Unlock()

	if ok == true {
		return content, nil
	} else {
		return nil, fmt.Errorf("humans.GetContent: No question asked for id \"%s\".", id)
	}

}

func (tool *Humans) MarshalJSON() ([]byte, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	return json.Marshal(tool.contents)

}

// NOTE: Await blocks until the human answered the question. This mirrors
// agents.Await deliberately: the planner model would otherwise poll in a hot
// loop and blow up its limited context window with identical "still waiting"
// tool messages. The lifecycle guarantees the done channel always closes
// (answer received -> Answer() -> close, or timeout -> watch_question ->
// close), so Await never blocks forever.
func (tool *Humans) Await(id string) (string, error) {

	tool.mutex.Lock()
	state, waiting := tool.states[id]
	tool.mutex.Unlock()

	if waiting == true {
		<-state.done
	}

	tool.mutex.Lock()
	question, ok := tool.contents[id]
	answer := ""

	if ok == true {
		answer = question.Answer
	}

	tool.mutex.Unlock()

	if ok == false {
		return "", fmt.Errorf("humans.Await: Question \"%s\" was never asked!", id)
	}

	if answer != "" {
		return fmt.Sprintf("humans.Await: Answer for question \"%s\"\n===%s\n===", id, answer), nil
	} else {
		return "", fmt.Errorf("humans.Await: Question \"%s\" was never answered!", id)
	}

}

