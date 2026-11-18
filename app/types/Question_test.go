package types

import "testing"

func TestNewQuestion(t *testing.T) {

	question := NewQuestion(map[string]any{
		"question": "Mascot?",
		"options":  []any{"Gnu", "Penguin"},
		"multiple": false,
	})

	if question.Type != "Choose" {
		t.Errorf("Expected %q, got %q", "Choose", question.Type)
	}

	if len(question.Options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(question.Options))
	}

}

func TestNewQuestion_Ask(t *testing.T) {

	question := NewQuestion(map[string]any{
		"question": "How old are you?",
	})

	if question.Type != "Ask" {
		t.Errorf("Expected %q, got %q", "Ask", question.Type)
	}

}

func TestNewQuestion_Permission(t *testing.T) {

	question := NewQuestion(map[string]any{
		"reason": "Need target access.",
		"capabilities": []any{
			map[string]any{"type": "files.Read", "scope": "/home/foo"},
		},
	})

	if question.Type != "Permission" {
		t.Errorf("Expected %q, got %q", "Permission", question.Type)
	}

	if len(question.Capabilities) != 1 {
		t.Fatalf("Expected 1 capability, got %d", len(question.Capabilities))
	}

	if question.Capabilities[0].Type != "files.Read" {
		t.Errorf("Expected %q, got %q", "files.Read", question.Capabilities[0].Type)
	}

}
