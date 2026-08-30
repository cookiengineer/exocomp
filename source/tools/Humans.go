package tools

import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"
import "encoding/json"
import "fmt"
import "strings"
import "sync"

var question_unique_id int64 = 0

type question_state struct {
	cancel context.CancelFunc
	done   chan struct{}
}

type Humans struct {
	contents map[string]*types.Question
	states   map[string]*question_state
	mutex    *sync.Mutex
}

func NewHumans() *humans {

	humans := &Humans{
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

	if ok == true {

		question.Answer = answer

		return nil

	} else {
		return fmt.Errorf("humans.Answer: Invalid question id.")
	}

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

	// TODO: How to create cancel context?
	state := &question_state{
		cancel: cancel,
		done:   make(chan struct{}),
	}

	tool.contents[question.ID] = question
	tool.states[question.ID] = state

	// TODO: go watch_answer(tool, state) ??? like in agents?

	return tool.Await(question.ID)

}

func (tool *Humans) Choice(text string, options []string, multiple bool) (string, error) {

	question_unique_id++

	question := &types.Question{
		ID:       fmt.Sprintf("question-%d", question_unique_id),
		Type:     "Choice",
		Question: text,
		Options:  []string{},
		Multiple: false,
		Answer:   "",
	}

	// TODO: How to create cancel context?
	state := &question_state{
		cancel: cancel,
		done:   make(chan struct{}),
	}

	tool.contents[question.ID] = question
	tool.states[question.ID] = state

	// TODO: go watch_answer(tool, state) ??? like in agents?

	return tool.Await(question.ID)

}

func (tool *Humans) GetContent(id string) (any, error) {

	content, ok := tool.contents[id]

	if ok == true {
		return content, nil
	} else {
		return nil, fmt.Errorf("requirements.GetContent: No question asked for id \"%s\".", id)
	}

}

func (tool *Humans) MarshalJSON() ([]byte, error) {
	return json.Marshal(tool.contents)
}

func (tool *Humans) Await(id string) (string, error) {

	question, ok := tool.contents[id]

	// TODO: What needs to happen is now a blocking channel
	// which is marked as done when the answer was written

	if ok == true {

		if question.Answer != "" {
			return fmt.Sprintf("humans.Await: Question \"%s\" was answered.", id), nil
		} else {
			return fmt.Sprintf("humans.Await: Question \"%s\" wasn't answered.", id), nil
		}

	} else {
		return "", fmt.Sprintf("humans.Await: Question \"%s\" was never asked!", id)
	}

}

