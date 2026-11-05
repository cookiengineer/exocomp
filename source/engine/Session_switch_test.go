package engine

import "exocomp/schemas"
import "exocomp/types"
import "testing"

func TestSession_SwitchAgent(t *testing.T) {

	session := newTestSession()
	session.Agent.Messages = append(session.Agent.Messages, &schemas.Message{
		Role:    "user",
		Content: "This message must not survive the switch.",
	})

	agent := &types.Agent{
		Name:        "Peanut Coder",
		Role:        "coder",
		Model:       "coder-model",
		Temperature: 0.5,
		Messages: []*schemas.Message{{
			Role:    "system",
			Content: "You are the coder profile.",
		}},
	}

	err0 := session.SwitchAgent(agent, "The previous conversation summary.")

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if session.Agent != agent {
		t.Errorf("Expected the switched agent to be installed")
	}

	if len(session.Agent.Messages) != 2 {
		t.Fatalf("Expected 2 messages after switch, got %d", len(session.Agent.Messages))
	}

	if session.Agent.Messages[0].Role != "system" {
		t.Errorf("Expected first message role %q, got %q", "system", session.Agent.Messages[0].Role)
	}

	if session.Agent.Messages[0].Content != "You are the coder profile." {
		t.Errorf("Expected the new system prompt, got %s", session.Agent.Messages[0].Content)
	}

	if session.Agent.Messages[1].Role != "assistant" {
		t.Errorf("Expected second message role %q, got %q", "assistant", session.Agent.Messages[1].Role)
	}

	if session.Agent.Messages[1].Content != "The previous conversation summary." {
		t.Errorf("Expected the history summary, got %s", session.Agent.Messages[1].Content)
	}

	if session.Config.Name != "Peanut Coder" {
		t.Errorf("Expected config name %q, got %q", "Peanut Coder", session.Config.Name)
	}

	if session.Config.Role != "coder" {
		t.Errorf("Expected config role %q, got %q", "coder", session.Config.Role)
	}

	if session.Config.Model != "coder-model" {
		t.Errorf("Expected config model %q, got %q", "coder-model", session.Config.Model)
	}

	if session.Config.Temperature != 0.5 {
		t.Errorf("Expected config temperature %v, got %v", 0.5, session.Config.Temperature)
	}

}

type switchToolStub struct{}

func (tool *switchToolStub) Name() string {
	return "agents"
}

func (tool *switchToolStub) Call(method string, arguments map[string]any) (string, error) {
	return "agents.Switch: Switched to agent profile.", nil
}

func (tool *switchToolStub) GetContent(id string) (any, error) {
	return nil, nil
}

func (tool *switchToolStub) GetContentIdentifiers() []string {
	return []string{}
}

func (tool *switchToolStub) HasMethod(method string) bool {
	return method == "Switch"
}

func (tool *switchToolStub) Schemas() []schemas.Tool {
	return []schemas.Tool{}
}

func TestSession_CallTool_SwitchDoesNotAppendToolMessage(t *testing.T) {

	session := newTestSession()
	session.SetTool(&switchToolStub{})

	err0 := session.CallTool("call_1", "agents", "Switch", map[string]any{
		"role": "coder",
	})

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if len(session.Agent.Messages) != 0 {
		t.Errorf("Expected no tool result to be appended, got %d messages", len(session.Agent.Messages))
	}

}
