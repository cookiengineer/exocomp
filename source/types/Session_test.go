package types

import "exocomp/schemas"
import "encoding/json"
import "fmt"
import "io"
import "net/http"
import net_url "net/url"
import "os"
import "strings"
import "sync"
import "testing"

func newTestSession() *Session {

	url, _ := net_url.Parse("http://localhost:11434/v1")
	config := NewConfig("Test Agent", "planner", "qwen3-coder:30b", "", 0.5, "/tmp/exocomp-test", "/tmp/exocomp-test", url, false)

	return &Session{
		Agent:    &Agent{Messages: make([]*schemas.Message, 0)},
		Config:   config,
		Console:  NewConsole(os.Stdout, os.Stderr, 0),
		Recovery: NewRecovery(config.Playground),
		Tools:    make([]*schemas.Tool, 0),
		Waiting:  false,
		client:   &http.Client{},
		mutex:    &sync.RWMutex{},
		tools:    make(map[string]Tool),
	}

}

type stubAgentsTool struct {
	callCount int
}

func (tool *stubAgentsTool) Call(method string, arguments map[string]interface{}) (string, error) {

	tool.callCount++

	if method == "Await" {
		return "", fmt.Errorf("agents.Await: Agent \"X\" is still working ...")
	}

	return "ok", nil

}

func (tool *stubAgentsTool) Get(id string) (any, error) {
	return nil, fmt.Errorf("stubAgentsTool.Get: nope")
}

type mockTransport struct {
	calls int
}

func (transport *mockTransport) RoundTrip(request *http.Request) (*http.Response, error) {

	transport.calls++

	body := `{"choices":[{"message":{"role":"assistant","content":"done"}}]}`

	return &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil

}

func TestSession_GetTool_NoDot(t *testing.T) {

	session := newTestSession()

	session.Tools = append(session.Tools, &schemas.Tool{
		Type:     "function",
		Function: schemas.ToolFunction{Name: "agentsList"},
	})

	result := session.GetTool("agentsList")

	if result != nil {
		t.Errorf("Expected %v to be nil", result)
	}

}

func TestSession_ReceiveChatResponse_NoAwaitPollingLoop(t *testing.T) {

	session := newTestSession()
	stub      := &stubAgentsTool{}
	transport := &mockTransport{}

	session.SetTool("agents", stub, []schemas.Tool{{
		Type:     "function",
		Function: schemas.ToolFunction{Name: "agents.Await"},
	}})

	session.client = &http.Client{Transport: transport}

	response := schemas.Message{
		Role:    "assistant",
		Content: "I will wait for the agent now.",
		ToolCalls: []schemas.ToolCall{{
			ID:   "call_1",
			Type: "function",
			Function: schemas.ToolCallFunction{
				Name:         "agents.Await",
				ArgumentsRaw: json.RawMessage(`{"name":"X"}`),
			},
		}},
	}

	err0 := session.ReceiveChatResponse(response)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if stub.callCount != 1 {
		t.Errorf("Expected the Await tool to be called exactly once, got %d", stub.callCount)
	}

	if transport.calls != 1 {
		t.Errorf("Expected exactly one follow-up inference, got %d", transport.calls)
	}

}
