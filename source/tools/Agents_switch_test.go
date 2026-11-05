package tools

import "exocomp/schemas"
import net_url "net/url"
import "path/filepath"
import "strings"
import "testing"

func TestAgents_Switch(t *testing.T) {

	playground := t.TempDir()
	sandbox := filepath.Join(playground, "work")
	url, _ := net_url.Parse("http://localhost:11434/v1")
	tool := NewAgents([]string{"Switch"}, playground, sandbox, "test-model", url, false)

	called := make(chan string, 1)

	tool.OnSwitch = func(role string) (string, error) {
		called <- role
		return "agents.Switch: Switched to agent profile.", nil
	}

	result, err0 := tool.Switch("Coder")

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.HasPrefix(result, "agents.Switch:") {
		t.Errorf("Expected a switch report, got %s", result)
	}

	select {
	case role := <-called:

		if role != "coder" {
			t.Errorf("Expected role %q, got %q", "coder", role)
		}

	default:
		t.Errorf("Expected OnSwitch hook to be called")
	}

}

func TestAgents_Switch_NoHook(t *testing.T) {

	playground := t.TempDir()
	sandbox := filepath.Join(playground, "work")
	url, _ := net_url.Parse("http://localhost:11434/v1")
	tool := NewAgents([]string{"Switch"}, playground, sandbox, "test-model", url, false)

	_, err0 := tool.Switch("coder")

	if err0 == nil {
		t.Fatalf("Expected a non-nil error")
	}

	if !strings.Contains(err0.Error(), "main session") {
		t.Errorf("Expected a main session error, got %v", err0)
	}

}

func TestAgents_Switch_UnknownRole(t *testing.T) {

	playground := t.TempDir()
	sandbox := filepath.Join(playground, "work")
	url, _ := net_url.Parse("http://localhost:11434/v1")
	tool := NewAgents([]string{"Switch"}, playground, sandbox, "test-model", url, false)

	tool.OnSwitch = func(role string) (string, error) {
		return "agents.Switch: Switched to agent profile.", nil
	}

	_, err0 := tool.Switch("does-not-exist")

	if err0 == nil {
		t.Fatalf("Expected a non-nil error")
	}

	if !strings.Contains(err0.Error(), "Unknown agent role") {
		t.Errorf("Expected an unknown role error, got %v", err0)
	}

}

func TestAgents_Call_Switch(t *testing.T) {

	playground := t.TempDir()
	sandbox := filepath.Join(playground, "work")
	url, _ := net_url.Parse("http://localhost:11434/v1")
	tool := NewAgents([]string{"Switch"}, playground, sandbox, "test-model", url, false)

	tool.OnSwitch = func(role string) (string, error) {
		return "agents.Switch: Switched to agent profile.", nil
	}

	result, err0 := tool.Call("Switch", map[string]interface{}{
		"role": "Coder",
	})

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.HasPrefix(result, "agents.Switch:") {
		t.Errorf("Expected a switch report, got %s", result)
	}

	_, err1 := tool.Call("Switch", map[string]interface{}{
		"role": 42,
	})

	if err1 == nil {
		t.Errorf("Expected a non-nil error")
	}

	if !strings.Contains(err1.Error(), "Invalid parameter \"role\"") {
		t.Errorf("Expected a role validation error, got %v", err1)
	}

}

func TestAgents_Summarize(t *testing.T) {

	tool := newTestAgents(t, "summarize")

	result, err0 := tool.Summarize([]*schemas.Message{{
		Role:    "user",
		Content: "Please implement the feature.",
	}})

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.Contains(result, "The agent implemented the feature successfully.") {
		t.Errorf("Expected summarized content, got %s", result)
	}

}
