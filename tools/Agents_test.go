//go:build agents

package tools

import net_url "net/url"
import "os"
import "path/filepath"
import "strings"
import "testing"
import "time"

import "fmt"

func TestAgents_Hire(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-agents-*")
	sandbox       := filepath.Join(playground, "agents")
	model         := "qwen3-coder:30b"
	url,        _ := net_url.Parse("http://localhost:11434/v1")
	tool          := NewAgents(playground, sandbox, model, url, true)

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
				"",
				"You MUST write the final result to the filesystem before you're finished.",
			}, "\n"),
		)

		done := make(chan bool, 1)

		go func() {

			ticker  := time.NewTicker(10 * time.Second)
			timeout := time.After(5 * time.Minute)

			for {
				select {
				case <-ticker.C:

					tool.Mutex.Lock()
					_, ok := tool.processes["Fibonacci Worker"]
					tool.Mutex.Unlock()

					if ok == false {
						done<-true
						return
					}

				case <-timeout:
					done<-false
					return
				}
			}

			ticker.Stop()

		}()

		finished := <-done

		if finished == true {

			fmt.Println("Agent finished!")

		} else {
			t.Errorf("Expected agent %s to finish", "Fibonacci Worker")
		}

		fmt.Println(finished)

		// TODO: goroutine that waits for tool.processes[agent_name]*os.Process to be nil
		// (which means the process is finished)

		// TODO: use tool.List() to show overview of working agents

		fmt.Println(result1)
		fmt.Println(err1)

	} else {
		t.Errorf("Expected %v to be not nil", tool)
	}

	// TODO
	// t.Cleanup(func() {

	// 	if t.Failed() == true {
	// 		t.Logf("Preserving folder %s for debugging.", playground)
	// 	} else {
	// 		os.RemoveAll(playground)
	// 	}

	// })

}
