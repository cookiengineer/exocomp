package chat

import app_types "exocomp-app/types"
import "exocomp-app/schemas"
import "encoding/json"
import "testing"

func newAssistantQuestion(id string, name string, arguments string) *schemas.Message {

	message := schemas.Message{
		Role: "assistant",
		ToolCalls: []schemas.ToolCall{{
			ID:   id,
			Type: "function",
			Function: schemas.ToolCallFunction{
				Name:         name,
				ArgumentsRaw: json.RawMessage(arguments),
			},
		}},
	}

	return &message

}

func TestGetUnansweredQuestions_Unanswered(t *testing.T) {

	agent := app_types.NewAgent()
	agent.Messages = append(agent.Messages, newAssistantQuestion("call_0", "humans.Ask", "{\"question\":\"How old are you?\"}"))

	questions := GetUnansweredQuestions(agent)

	if len(questions) != 1 {
		t.Fatalf("Expected 1 unanswered question, got %d", len(questions))
	}

	if questions[0].Question != "How old are you?" {
		t.Errorf("Expected %q, got %q", "How old are you?", questions[0].Question)
	}

	if questions[0].Type != "Ask" {
		t.Errorf("Expected %q, got %q", "Ask", questions[0].Type)
	}

}

func TestGetUnansweredQuestions_Answered(t *testing.T) {

	agent := app_types.NewAgent()
	agent.Messages = append(agent.Messages, newAssistantQuestion("call_0", "humans.Ask", "{\"question\":\"How old are you?\"}"))
	agent.Messages = append(agent.Messages, &schemas.Message{
		Role:       "tool",
		ToolName:   "humans.Ask",
		ToolCallID: "call_0",
		Content:    "Answer for Question\n===\n42",
	})

	questions := GetUnansweredQuestions(agent)

	if len(questions) != 0 {
		t.Errorf("Expected 0 unanswered questions, got %d", len(questions))
	}

}

func TestGetUnansweredQuestions_Choice(t *testing.T) {

	agent := app_types.NewAgent()
	agent.Messages = append(agent.Messages, newAssistantQuestion("call_0", "humans.Choose", "{\"question\":\"Mascot?\",\"options\":[\"Gnu\",\"Penguin\"],\"multiple\":false}"))

	questions := GetUnansweredQuestions(agent)

	if len(questions) != 1 {
		t.Fatalf("Expected 1 unanswered question, got %d", len(questions))
	}

	if questions[0].Type != "Choose" {
		t.Errorf("Expected %q, got %q", "Choose", questions[0].Type)
	}

	if len(questions[0].Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(questions[0].Options))
	}

}
