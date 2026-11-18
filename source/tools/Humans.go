package tools

import "exocomp/schemas"
import "exocomp/types"
import utils_capabilities "exocomp/utils/capabilities"
import utils_fmt "exocomp/utils/fmt"
import "context"
import "encoding/json"
import "fmt"
import "slices"
import "sort"
import "strings"
import "sync"
import "time"

type questionState struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

type Humans struct {
	Methods      []string
	AgentName    string
	Capabilities *utils_capabilities.Set
	Timeout      time.Duration
	contents     map[string]*types.Question
	states       map[string]*questionState
	mutex        *sync.RWMutex
}

func NewHumans(methods []string, playground string, sandbox string) *Humans {

	return NewHumansWithCapabilities(methods, playground, sandbox, "", nil)

}

func NewHumansWithCapabilities(
	methods []string,
	playground string,
	sandbox string,
	agent_name string,
	capabilities *utils_capabilities.Set,
) *Humans {

	humans := &Humans{
		Methods:      methods,
		AgentName:    agent_name,
		Capabilities: capabilities,
		Timeout:      15 * time.Minute,
		contents:     make(map[string]*types.Question),
		states:       make(map[string]*questionState),
		mutex:        &sync.RWMutex{},
	}

	if humans.Capabilities == nil {
		humans.Capabilities = utils_capabilities.NewSet()
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
			raw, ok2 := arguments["options"]
			multiple, _ := arguments["multiple"].(bool)

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

		} else if method == "AskPermission" {

			reason, ok1 := arguments["reason"].(string)
			raw_capabilities, ok2 := arguments["capabilities"]

			if ok1 == true && ok2 == true {

				capabilities, err0 := utils_capabilities.Parse(raw_capabilities)

				if err0 == nil {
					return tool.AskPermission(utils_fmt.FormatSingleLine(reason), capabilities)
				} else {
					return "", fmt.Errorf("humans.%s: %s", method, err0.Error())
				}

			} else if ok1 == true && ok2 == false {
				return "", fmt.Errorf("humans.%s: %s", method, "Invalid parameter \"capabilities\" is not an array.")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("humans.%s: %s", method, "Invalid parameter \"reason\" is not a string.")
			} else {
				return "", fmt.Errorf("humans.%s: Invalid parameters.", method)
			}

		} else if method == "Answer" {

			question, ok1 := arguments["question"].(string)
			answer, ok2 := arguments["answer"].(string)

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

func (tool *Humans) AskPermission(reason string, capabilities []types.Capability) (string, error) {

	if len(capabilities) > 32 {
		return "", fmt.Errorf("humans.AskPermission: Too many capabilities (maximum 32).")
	}

	missing, err0 := tool.Capabilities.Missing(capabilities)

	if err0 == nil {

		if len(missing) > 0 {

			question_text := strings.TrimSpace(reason)

			if question_text == "" {
				question_text = "Permission request"
			}

			question := &types.Question{
				Type:         "Permission",
				Agent:        tool.AgentName,
				Question:     question_text,
				Options:      []string{},
				Multiple:     true,
				Capabilities: missing,
				Answer:       "",
			}

			ctx, cancel := context.WithTimeout(context.Background(), tool.Timeout)
			state := &questionState{
				ctx:    ctx,
				cancel: cancel,
				done:   make(chan struct{}),
			}

			tool.mutex.Lock()

			_, pending := tool.states[question.Question]

			if pending == true {
				tool.mutex.Unlock()
				cancel()
				return "", fmt.Errorf("humans.AskPermission: Permission request \"%s\" is already pending.", question.Question)
			}

			tool.contents[question.Question] = question
			tool.states[question.Question] = state
			tool.mutex.Unlock()

			go watchQuestion(tool, question.Question, state)

			_, err1 := tool.Await(question.Question)

			if err1 == nil {

				tool.mutex.RLock()
				answer_text := question.Answer
				tool.mutex.RUnlock()

				answer, err2 := utils_capabilities.ParseAnswer(answer_text)

				if err2 == nil {

					requested := make(map[string]types.Capability, len(missing))

					for _, capability := range missing {
						requested[capability.Key()] = capability
					}

					granted := make([]types.Capability, 0)
					denied := make([]types.Capability, 0)

					granted_keys := make(map[string]bool, len(answer))

					for _, capability := range answer {

						normalized, err3 := utils_capabilities.Normalize(capability)

						if err3 == nil {

							if requested_capability, ok := requested[normalized.Key()]; ok == true && granted_keys[normalized.Key()] == false {
								granted = append(granted, requested_capability)
								granted_keys[normalized.Key()] = true
							}

						}

					}

					for _, capability := range missing {

						if granted_keys[capability.Key()] == false {
							denied = append(denied, capability)
						}

					}

					if len(granted) > 0 {

						err4 := tool.Capabilities.Grant(granted)

						if err4 != nil {
							return "", fmt.Errorf("humans.AskPermission: %s", err4.Error())
						}

					}

					return permissionSummary(granted, denied), nil

				} else {
					return "", fmt.Errorf("humans.AskPermission: Invalid permission answer.")
				}

			} else {

				tool.mutex.Lock()

				if existing, ok := tool.contents[question.Question]; ok == true && existing.Answer == "" {
					delete(tool.contents, question.Question)
				}

				tool.mutex.Unlock()

				return "", err1
			}

		} else {
			return "humans.AskPermission: All requested capabilities are already granted.", nil
		}

	} else {
		return "", fmt.Errorf("humans.AskPermission: %s", err0.Error())
	}

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
