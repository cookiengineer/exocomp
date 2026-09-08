package tools

import "exocomp/types"
import "os"
import "path/filepath"
import "strings"
import "testing"

func TestAgents_List(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-agents-*")
	sandbox       := filepath.Join(playground, "agents")
	tool          := NewAgents([]string{"List", "Roles"}, playground, sandbox, "huihui_ai/Qwen3.6-abliterated:35b", nil, false)

	if tool != nil {

		result1, err1 := tool.List()

		if result1 != "" {
			t.Errorf("Expected %s to be empty", result1)
		}

		if err1 == nil {
			t.Errorf("Expected %v to be not nil", err1)
		}

		agent := &types.Agent{
			Name:   "Coder",
			Role:   "coder",
			Status: "working",
		}

		if tool.SetAgent(agent) == false {
			t.Fatalf("Expected agent to be registered")
		}

		result2, err2 := tool.List()

		if err2 != nil {
			t.Errorf("Expected %v to be nil", err2)
		}

		if strings.Contains(result2, "agents.List: 1 agents were working for us.") == false {
			t.Errorf("Expected one agent to be listed, got %s", result2)
		}

		if strings.Contains(result2, "Name: \"Coder\", Type: coder, Status: working") == false {
			t.Errorf("Expected agent details, got %s", result2)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestAgents_Roles(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-agents-*")
	sandbox       := filepath.Join(playground, "agents")
	tool          := NewAgents([]string{"List", "Roles"}, playground, sandbox, "huihui_ai/Qwen3.6-abliterated:35b", nil, false)

	if tool != nil {

		result, err := tool.Roles()

		if err != nil {
			t.Errorf("Expected %v to be nil", err)
		}

		if strings.Contains(result, "agents.Roles:") == false {
			t.Errorf("Expected roles header, got %s", result)
		}

		if strings.Contains(result, "Role: \"coder\"") == false {
			t.Errorf("Expected coder role to be available, got %s", result)
		}

		if strings.Contains(result, "Role: \"planner\"") == true {
			t.Errorf("Expected planner role to be excluded, got %s", result)
		}

	} else {
		t.Errorf("Expected tool to be not nil")
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}

func TestAgents_Inquire(t *testing.T) {

	tool := newTestAgents(t, "quit")
	name := hireTestAgent(t, tool, "quit")

	_, err0 := tool.Await(name)

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	t.Setenv("EXOCOMP_FAKE_SCENARIO", "summarize")

	result, err1 := tool.Inquire(name)

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	if strings.Contains(result, "agents.Inquire: Summary of already finished agent") == false {
		t.Errorf("Expected summary of finished agent, got %s", result)
	}

	if strings.Contains(result, "The agent implemented the feature successfully.") == false {
		t.Errorf("Expected summarized content, got %s", result)
	}

}
