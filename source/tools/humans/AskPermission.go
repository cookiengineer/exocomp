package humans

import "exocomp/types"
import utils_capabilities "exocomp/utils/capabilities"
import "context"
import "fmt"
import "strings"

func (tool *Humans) AskPermission(reason string, capabilities []types.Capability) (string, error) {

	if len(capabilities) <= 32 {

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

	} else {
		return "", fmt.Errorf("humans.AskPermission: Too many capabilities (maximum 32).")
	}
}
