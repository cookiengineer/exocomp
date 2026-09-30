package chat

import app_types "exocomp-app/types"
import "strings"

func GetUnansweredQuestions(agent *app_types.Agent) []*app_types.Question {

	questions := make(map[string]*app_types.Question)

	if agent == nil {
		return make([]*app_types.Question, 0)
	}

	for _, message := range agent.Messages {

		if message == nil {
			continue
		}

		if message.Role == "assistant" {

			for t := 0; t < len(message.ToolCalls); t++ {

				tool_call := &message.ToolCalls[t]

				if tool_call.Function.Name == "humans.Ask" || tool_call.Function.Name == "humans.Choose" {

					arguments, err := tool_call.GetArguments()

					if err == nil {
						questions[tool_call.ID] = app_types.NewQuestion(arguments)
					}

				}

			}

		} else if message.Role == "tool" {

			if message.ToolName == "humans.Ask" || message.ToolName == "humans.Choose" {

				question, ok := questions[message.ToolCallID]

				if ok == true && question != nil {

					content := strings.TrimSpace(message.Content)

					if strings.Contains(content, "Answer for Question") == true && strings.Contains(content, "===") == true {
						question.Answer = strings.Split(content, "===")[1]
					}

				}

			}

		}

	}

	unanswered := make([]*app_types.Question, 0)

	for _, question := range questions {

		if question.Answer == "" {
			unanswered = append(unanswered, question)
		}

	}

	return unanswered

}
