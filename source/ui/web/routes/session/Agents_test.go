package session

import "exocomp/engine"
import "exocomp/schemas"
import "exocomp/types"
import "encoding/json"
import "net/http"
import "net/http/httptest"
import "testing"

func newRouteTestSession(t *testing.T) *engine.Session {

	t.Helper()

	config := types.NewConfig(
		"Planner",
		"planner",
		"test-model",
		"",
		0.0,
		t.TempDir(),
		t.TempDir(),
		nil,
		false,
	)

	session := engine.NewSession(&types.Agent{
		Name:         "Planner",
		Role:         "planner",
		Model:        "test-model",
		Sandbox:      config.Sandbox,
		AllowedTools: []string{"agents.List"},
	}, config)

	session.InstallTools()

	session.InstallAgents([]*types.Agent{
		{
			Name:    "Coder",
			Role:    "coder",
			Model:   "test-model",
			Sandbox: config.Sandbox,
			Status:  "working",
			Messages: []*schemas.Message{
				{Role: "user", Content: "hello from human"},
			},
		},
	})

	return session

}

func TestAgentsRouteReturnsPlannerAndHiredAgents(t *testing.T) {

	session := newRouteTestSession(t)
	request := httptest.NewRequest(http.MethodGet, "/api/session/agents", nil)
	response := httptest.NewRecorder()

	Agents(session, request, response)

	if response.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, response.Code)
	}

	agents := make([]*types.Agent, 0)
	err := json.Unmarshal(response.Body.Bytes(), &agents)

	if err != nil {
		t.Fatalf("Expected valid agent JSON, got %v", err)
	}

	if len(agents) != 2 {
		t.Fatalf("Expected two agents, got %d", len(agents))
	}

	if agents[0].Name != "Planner" || agents[1].Name != "Coder" {
		t.Errorf("Expected planner first and hired agent second")
	}

	if len(agents[1].Messages) != 1 || agents[1].Messages[0].Content != "hello from human" {
		t.Errorf("Expected hired agent messages in the response")
	}

}

func TestMessagesRouteReturnsSelectedAgent(t *testing.T) {

	session := newRouteTestSession(t)
	request := httptest.NewRequest(http.MethodGet, "/api/session/messages?name=Coder", nil)
	response := httptest.NewRecorder()

	Messages(session, request, response)

	if response.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, response.Code)
	}

	messages := make([]*schemas.Message, 0)
	err := json.Unmarshal(response.Body.Bytes(), &messages)

	if err != nil {
		t.Fatalf("Expected valid message JSON, got %v", err)
	}

	if len(messages) != 1 || messages[0].Content != "hello from human" {
		t.Errorf("Expected the selected agent messages")
	}

}

func TestAgentRouteReturnsSelectedAgent(t *testing.T) {

	session := newRouteTestSession(t)
	request := httptest.NewRequest(http.MethodGet, "/api/session/agent?name=Coder", nil)
	response := httptest.NewRecorder()

	Agent(session, request, response)

	if response.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, response.Code)
	}

	agent := &types.Agent{}
	err := json.Unmarshal(response.Body.Bytes(), agent)

	if err != nil {
		t.Fatalf("Expected valid agent JSON, got %v", err)
	}

	if agent.Name != "Coder" {
		t.Errorf("Expected the selected agent, got %s", agent.Name)
	}

	if len(agent.Messages) != 1 {
		t.Errorf("Expected the selected agent messages")
	}

}

func TestAgentConfigRouteReturnsSelectedAgentConfig(t *testing.T) {

	session := newRouteTestSession(t)
	request := httptest.NewRequest(http.MethodGet, "/api/session/config/Coder", nil)
	request.SetPathValue("name", "Coder")
	response := httptest.NewRecorder()

	AgentConfig(session, request, response)

	if response.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, response.Code)
	}

	config := &types.Config{}
	err := json.Unmarshal(response.Body.Bytes(), config)

	if err != nil {
		t.Fatalf("Expected valid config JSON, got %v", err)
	}

	if config.Name != "Coder" || config.Role != "coder" {
		t.Errorf("Expected the selected agent config")
	}

}
