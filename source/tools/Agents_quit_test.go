package tools

import net_url "net/url"
import "strings"
import "testing"
import "time"

func TestAgents_Quit_CallsHook(t *testing.T) {

	playground := t.TempDir()
	url, _     := net_url.Parse("http://localhost:11434/v1")
	tool       := NewAgents(playground, sandbox, "huihui_ai/Qwen3.6-abliterated:35b", url, false)

	type quit_call struct {
		report  string
		success bool
	}

	called := make(chan quit_call, 1)

	tool.OnQuit = func(report string, success bool) {
		called <- quit_call{report: report, success: success}
	}

	report, err0 := tool.Quit("my work is done, here is the summary")

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.HasPrefix(report, "agents.Quit: Work Report\n") {
		t.Errorf("Expected report to be a work report, got %s", report)
	}

	select {
	case call := <-called:

		if call.success != true {
			t.Errorf("Expected success to be true")
		}

		if call.report != "my work is done, here is the summary" {
			t.Errorf("Expected report %q, got %q", "my work is done, here is the summary", call.report)
		}

	case <-time.After(time.Second):
		t.Errorf("Expected OnQuit hook to be called")
	}

}

func TestAgents_Quit_FailureHook(t *testing.T) {

	playground := t.TempDir()
	url, _     := net_url.Parse("http://localhost:11434/v1")
	tool       := NewAgents(playground, sandbox, "huihui_ai/Qwen3.6-abliterated:35b", url, false)

	called := make(chan bool, 1)

	tool.OnQuit = func(report string, success bool) {
		called <- success
	}

	_, err0 := tool.Quit("I give up, this cannot be done")

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	select {
	case success := <-called:

		if success != false {
			t.Errorf("Expected success to be false")
		}

	case <-time.After(time.Second):
		t.Errorf("Expected OnQuit hook to be called")
	}

}
