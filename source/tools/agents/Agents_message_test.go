package agents

import "strings"
import "testing"

func TestAgents_SendMessage(t *testing.T) {

	tool := newTestAgents(t, "input")
	name := hireTestAgent(t, tool, "input")

	result1, err1 := tool.SendMessage(name, "hello from human")

	if err1 != nil {
		t.Errorf("Expected %v to be nil", err1)
	}

	if strings.Contains(result1, "Message sent to agent") == false {
		t.Errorf("Expected a send confirmation, got %s", result1)
	}

	_, err2 := tool.Await(name)

	if err2 != nil {
		t.Fatalf("Expected %v to be nil", err2)
	}

	agent := tool.GetAgentSnapshot(name)

	if agent == nil {
		t.Fatalf("Expected agent %s to exist", name)
	}

	found_user := false
	found_assistant := false

	for _, message := range agent.Messages {

		if message.Content == "hello from human" {
			found_user = true
		}

		if message.Content == "ack: hello from human" {
			found_assistant = true
		}

	}

	if found_user == false {
		t.Errorf("Expected the human message in the agent history")
	}

	if found_assistant == false {
		t.Errorf("Expected the agent reply in the agent history")
	}

	_, err3 := tool.SendMessage(name, "again")

	if err3 == nil {
		t.Errorf("Expected sending to a finished agent to fail")
	}

}

func TestAgents_SendMessage_UnknownAgent(t *testing.T) {

	tool := newTestAgents(t, "input")

	_, err := tool.SendMessage("Missing Agent", "hello")

	if err == nil {
		t.Errorf("Expected an error for an unknown agent")
	}

}
