package tools

import "exocomp/schemas"
import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"
import "context"
import "encoding/json"
import "fmt"
import "slices"
import "sort"
import "strings"
import "sync"
import "time"

type question_state struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

func watch_question(tool *Humans, question string, state *question_state) {

	<-state.ctx.Done()

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	_, still_waiting := tool.states[question]

	if still_waiting == true {
		delete(tool.states, question)
		close(state.done)
	}

}

type Humans struct {
	Methods  []string
	Timeout  time.Duration
	contents map[string]*types.Question
	states   map[string]*question_state
	mutex    *sync.RWMutex
}

func NewHumans(methods []string, playground string, sandbox string) *Humans {

	humans := &Humans{
		Methods:  methods,
		Timeout:  15 * time.Minute,
		contents: make(map[string]*types.Question),
		states:   make(map[string]*question_state),
		mutex:    &sync.RWMutex{},
	}

	return humans

}

func (tool *Humans) Name() string {
	return "humans"
}

func (tool *Humans) Call(method string, arguments map[string]interface{}) (string, error) {

	if tool.HasMethod(method) == true {

		if method == "Ask" {

			question, ok1 := arguments["question"].(string)

			if ok1 == true {
				return tool.Ask(utils_fmt.FormatSingleLine(question))
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

				return tool.Choose(utils_fmt.FormatSingleLine(question), options, multiple)

			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("humans.%s: %s", method, "Invalid parameter \"options\" is not an array of strings.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("humans.%s: %s", method, "Invalid parameter \"question\" is not a string.")
			} else {
				return "", fmt.Errorf("humans.%s: Invalid parameters.", method)
			}

		} else if method == "Answer" {

			question, ok1 := arguments["question"].(string)
			answer,   ok2 := arguments["answer"].(string)

			if ok1 == true && ok2 == true {
				return tool.Answer(utils_fmt.FormatSingleLine(question), utils_fmt.FormatMultiLine(answer))
			} else {
				return "", fmt.Errorf("humans.%s: %s", method, "Invalid parameters \"id\" and \"answer\" are not strings.")
			}

		} else {
			return "", fmt.Errorf("humans.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("humans.%s: Method not allowed.", method)
	}

}

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

func (tool *Humans) Ask(text string) (string, error) {

	question := &types.Question{
		Type:     "Ask",
		Question: strings.TrimSpace(text),
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
	tool.contents[question.Question] = question
	tool.states[question.Question]  = state
	tool.mutex.Unlock()

	go watch_question(tool, question.Question, state)

	return tool.Await(question.Question)

}

func (tool *Humans) Choose(text string, options []string, multiple bool) (string, error) {

	question := &types.Question{
		Type:     "Choose",
		Question: strings.TrimSpace(text),
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
	tool.contents[question.Question] = question
	tool.states[question.Question] = state
	tool.mutex.Unlock()

	go watch_question(tool, question.Question, state)

	return tool.Await(question.Question)

}

func (tool *Humans) GetContent(question string) (any, error) {

	question = strings.TrimSpace(question)

	tool.mutex.RLock()
	content, ok := tool.contents[question]
	tool.mutex.RUnlock()

	if ok == true {
		return content, nil
	} else {
		return nil, fmt.Errorf("humans.GetContent: Question \"%s\" was never asked!", question)
	}

}

func (tool *Humans) GetContentIdentifiers() []string {

	result := make([]string, 0)

	tool.mutex.RLock()

	for id, _ := range tool.contents {
		result = append(result, id)
	}

	tool.mutex.RUnlock()

	sort.Strings(result)

	return result

}

func (tool *Humans) HasMethod(method string) bool {

	// NOTE: Special case, required for UI/UX integration
	if method == "Answer" {
		return true
	} else {
		return slices.Contains(tool.Methods, method) == true
	}

}

func (tool *Humans) MarshalJSON() ([]byte, error) {

	tool.mutex.Lock()
	defer tool.mutex.Unlock()

	return json.Marshal(tool.contents)

}

func (tool *Humans) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range HumansSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}



// NOTE: Await blocks until the human answered the question. This mirrors
// agents.Await deliberately: the planner model would otherwise poll in a hot
// loop and blow up its limited context window with identical "still waiting"
// tool messages. The lifecycle guarantees the done channel always closes
// (answer received -> Answer() -> close, or timeout -> watch_question ->
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
			return fmt.Sprintf("humans.%s: Answer for Question \"%s\" is:\n===\n%s\n===", question.Type, question.Question, question.Answer), nil
		} else {
			return "", fmt.Errorf("humans.%s: Question \"%s\" was never answered!", question.Type, question.Question)
		}

	} else {
		return "", fmt.Errorf("humans.Await: Question \"%s\" was never asked!", reference)
	}

}

