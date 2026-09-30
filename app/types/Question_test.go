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
