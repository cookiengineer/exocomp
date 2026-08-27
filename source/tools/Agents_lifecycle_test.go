package tools

import net_url "net/url"
import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "runtime"
import "strings"
import "sync"
import "testing"
import "time"

var (
	fakeAgentOnce sync.Once
	fakeAgentPath string
	fakeAgentErr  error
)

func buildFakeAgent(t *testing.T) string {

	fakeAgentOnce.Do(func() {

		_, filename, _, ok := runtime.Caller(0)

		if ok == false {
			fakeAgentErr = fmt.Errorf("Cannot locate test file path")
			return
		}

		tools_dir := filepath.Dir(filename)
		output    := filepath.Join(os.TempDir(), "exocomp-fake-agent")

		cmd := exec.Command("go", "build", "-o", output, "./testdata/fakeagent")
		cmd.Dir = tools_dir

		data, err := cmd.CombinedOutput()

		if err != nil {
			fakeAgentErr = fmt.Errorf("Cannot build fake agent: %v: %s", err, strings.TrimSpace(string(data)))
			return
		}

		fakeAgentPath = output

	})

	if fakeAgentErr != nil {
		t.Fatalf("%v", fakeAgentErr)
	}

	return fakeAgentPath

}

func newTestAgents(t *testing.T, scenario string) *Agents {

	playground := t.TempDir()
	sandbox    := filepath.Join(playground, "work")
	url, _     := net_url.Parse("http://localhost:11434/v1")

	t.Setenv("EXOCOMP_AGENT", buildFakeAgent(t))
	t.Setenv("EXOCOMP_FAKE_SCENARIO", scenario)

	return NewAgents(playground, sandbox, "huihui_ai/Qwen3.6-abliterated:35b", url, false)

}

func hireTestAgent(t *testing.T, tool *Agents, scenario string) string {

	result, err0 := tool.Hire("coder", "Please implement the feature.", ".")

	if err0 != nil {
		t.Fatalf("Expected %v to be nil", err0)
	}

	if !strings.HasPrefix(result, "agents.Hire: Agent \"") {
		t.Fatalf("Expected %s to be a hire report", result)
	}

	names := tool.GetNames()

	if len(names) != 1 {
		t.Fatalf("Expected 1 hired agent, got %d", len(names))
	}

	return names[0]

}

func TestAgents_HireAndAwait_WorkReport(t *testing.T) {

	tool := newTestAgents(t, "quit")
	name := hireTestAgent(t, tool, "quit")

	report, err0 := tool.Await(name)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.Contains(report, "my work is done: implemented the feature") {
		t.Errorf("Expected work report to contain the quit message, got %s", report)
	}

	agent := tool.GetAgent(name)

	if agent == nil {
		t.Fatalf("Expected agent %s to exist", name)
	}

	if agent.Status != "finished" {
		t.Errorf("Expected agent status to be %q, got %q", "finished", agent.Status)
	}

}

func TestAgents_HireAndAwait_LargeMessage(t *testing.T) {

	tool := newTestAgents(t, "large")
	name := hireTestAgent(t, tool, "large")

	report, err0 := tool.Await(name)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.Contains(report, "my work is done: implemented the feature") {
		t.Errorf("Expected work report to contain the quit message, got %s", report)
	}

	agent := tool.GetAgent(name)

	if agent == nil {
		t.Fatalf("Expected agent %s to exist", name)
	}

	large_seen := false

	for _, message := range agent.Messages {

		if message.Role == "tool" && strings.HasPrefix(message.Content, "files.Read: big.go") {
			large_seen = true
			break
		}

	}

	if large_seen == false {
		t.Errorf("Expected the >64KiB message to be captured by the reader")
	}

}

func TestAgents_Await_IdleTimeout(t *testing.T) {

	tool := newTestAgents(t, "hang")
	tool.IdleTimeout = 200 * time.Millisecond
	tool.Timeout     = 10 * time.Second

	name := hireTestAgent(t, tool, "hang")

	start := time.Now()
	_, err0 := tool.Await(name)
	elapsed := time.Since(start)

	if err0 == nil {
		t.Errorf("Expected a non-nil error for a hanging agent")
	}

	if !strings.Contains(err0.Error(), "never finished with a work report") {
		t.Errorf("Expected never finished error, got %v", err0)
	}

	if elapsed > 10*time.Second {
		t.Errorf("Expected Await to return within ~10s, took %s", elapsed)
	}

	agent := tool.GetAgent(name)

	if agent != nil && agent.Status != "failed" {
		t.Errorf("Expected agent status to be %q, got %q", "failed", agent.Status)
	}

}

func TestAgents_Fire(t *testing.T) {

	tool := newTestAgents(t, "hang")
	name := hireTestAgent(t, tool, "hang")

	result, err0 := tool.Fire(name)

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.Contains(result, "fired") {
		t.Errorf("Expected fire report, got %s", result)
	}

	agent := tool.GetAgent(name)

	if agent == nil {
		t.Fatalf("Expected agent %s to exist", name)
	}

	if agent.Status != "fired" {
		t.Errorf("Expected agent status to be %q, got %q", "fired", agent.Status)
	}

}

func TestAgents_Await_NeverFinished(t *testing.T) {

	tool := newTestAgents(t, "noquit")
	name := hireTestAgent(t, tool, "noquit")

	_, err0 := tool.Await(name)

	if err0 == nil {
		t.Errorf("Expected a non-nil error for an agent without a work report")
	}

	if !strings.Contains(err0.Error(), "never finished with a work report") {
		t.Errorf("Expected never finished error, got %v", err0)
	}

	agent := tool.GetAgent(name)

	if agent != nil && agent.Status != "failed" {
		t.Errorf("Expected agent status to be %q, got %q", "failed", agent.Status)
	}

}
