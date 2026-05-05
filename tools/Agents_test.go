//go:build agents

package tools

import net_url "net/url"
import "os"
import "path/filepath"
import "strings"
import "testing"

import "fmt"

func TestAgents_Hire(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-agents-*")
	sandbox       := filepath.Join(playground, "agents")
	model         := "qwen3-coder:30b"
	url,        _ := net_url.Parse("http://localhost:11434/v1")
	tool          := NewAgents(playground, sandbox, model, url)

	if tool != nil {

		result1, err1 := tool.Hire(
			"Fibonacci Worker",
			"coder",
			"./fibonacci",
			strings.Join([]string{
				"Can you write a main.go for me that implements the fibonacci sequence?",
				"The first parameter should be the sequence/step parameter.",
				"",
				"\"program 2\" should return 1",
				"\"program 3\" should return 1",
				"\"program 4\" should return 2",
				"\"program 5\" should return 3",
				"",
				"... and so on",
			}, "\n"),
		)

		// TODO: use tool.List() to show overview of working agents

		fmt.Println(result1)
		fmt.Println(err1)

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	t.Cleanup(func() {

		if t.Failed() == true {
			t.Logf("Preserving folder %s for debugging.", playground)
		} else {
			os.RemoveAll(playground)
		}

	})

}
