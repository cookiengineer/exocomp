package engine

import "exocomp/schemas"
import "exocomp/types"
import "encoding/json"
import "fmt"
import "io"
import "math"
import "net/http"
import net_url "net/url"
import "os"
import "strings"
import "sync"
import "testing"

func newTestSession() *Session {

	url, _ := net_url.Parse("http://localhost:11434/v1")
	config := types.NewConfig("Test Agent", "planner", "test-model", "", 0.0, "/tmp/exocomp-test", "/tmp/exocomp-test", url, false)

	return &Session{
		Agent: &types.Agent{
			Model:    "test-model",
			Messages: make([]*schemas.Message, 0),
		},
		Config:   config,
		Console:  types.NewConsole(os.Stdout, os.Stderr, 0),
		Recovery: NewRecovery(config.Playground),
		Waiting:  false,
		client:   &http.Client{},
		mutex:    &sync.RWMutex{},
		adapters: make(map[string]types.Adapter),
		tools:    make(map[string]types.Tool),
	}

}

type stubAgentsTool struct {
	callCount int
}

func (tool *stubAgentsTool) Name() string {
	return "agents"
}

func (tool *stubAgentsTool) Call(method string, arguments map[string]any) (string, error) {

	tool.callCount++

	if method == "Await" {
		return "", fmt.Errorf("agents.Await: Agent \"X\" is still working ...")
	}

	return "ok", nil

}

func (tool *stubAgentsTool) GetContent(id string) (any, error) {
	return nil, fmt.Errorf("stubAgentsTool.GetContent: nope")
}

func (tool *stubAgentsTool) GetContentIdentifiers() []string {
	return []string{}
}

func (tool *stubAgentsTool) HasMethod(method string) bool {
	return method == "Await"
}

func (tool *stubAgentsTool) Schemas() []schemas.Tool {
	return []schemas.Tool{}
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
	session.SetTool(&stubAgentsTool{})

	result := session.GetTool("agentsList")

	if result != nil {
		t.Errorf("Expected %v to be nil", result)
	}

}

func TestSession_ReceiveChatResponse_NoAwaitPollingLoop(t *testing.T) {

	session := newTestSession()
	stub      := &stubAgentsTool{}
	transport := &mockTransport{}

	session.SetTool(stub)
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

type mockUsageTransport struct {
	body string
}

func (transport *mockUsageTransport) RoundTrip(request *http.Request) (*http.Response, error) {

	return &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(transport.body)),
	}, nil

}

func newDeepSeekTestSession() *Session {

	session := newTestSession()

	url, _ := net_url.Parse("https://api.deepseek.com")

	session.Config = &types.Config{
		URL: url,
		Providers: map[string]types.Provider{
			"deepseek-v4-pro": {
				URL:   url,
				Alias: "deepseek-v4-pro",
				Pricing: types.Pricing{
					InputPrice:       1.32,
					OutputPrice:      3.96,
					CachedInputPrice: 0.044,
				},
			},
		},
	}

	session.Agent.Model = "deepseek-v4-pro"

	return session

}

func TestSession_InferChatCompletions_TracksTokensAndCost(t *testing.T) {

	session := newDeepSeekTestSession()

	body := `{"choices":[{"message":{"role":"assistant","content":"done","reasoning_content":"thinking..."}}],"usage":{"prompt_tokens":100,"completion_tokens":50,"total_tokens":150,"prompt_tokens_details":{"cached_tokens":20},"completion_tokens_details":{"reasoning_tokens":40}}}`

	session.client = &http.Client{
		Transport: &mockUsageTransport{
			body: body,
		},
	}

	err0 := session.infer_chat_completions()

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	usage := session.Agent.ContextUsage

	if usage.Tokens != 100 {
		t.Errorf("Expected Tokens %d to be %d", usage.Tokens, 100)
	}

	if usage.PromptTokens != 100 {
		t.Errorf("Expected PromptTokens %d to be %d", usage.PromptTokens, 100)
	}

	if usage.CompletionTokens != 50 {
		t.Errorf("Expected CompletionTokens %d to be %d", usage.CompletionTokens, 50)
	}

	if usage.TotalTokens != 150 {
		t.Errorf("Expected TotalTokens %d to be %d", usage.TotalTokens, 150)
	}

	expected := (float64(20)*0.044 + float64(80)*1.32 + float64(50)*3.96) / 1000000.0

	if math.Abs(usage.Cost-expected) > 0.000000001 {
		t.Errorf("Expected Cost %v to be %v", usage.Cost, expected)
	}

	last := session.GetLastMessage()

	if last == nil {
		t.Fatalf("Expected a last message to be appended")
	}

	if last.ReasoningContent != "thinking..." {
		t.Errorf("Expected ReasoningContent %q to be %q", last.ReasoningContent, "thinking...")
	}

}
