//go:build agents

package tools

import net_url "net/url"
import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "strings"
import "testing"
import "time"

func waitForAgent(tool *Agents, name string) bool {

	done := make(chan bool, 1)

	go func() {

		ticker  := time.NewTicker(10 * time.Second)
		timeout := time.After(5 * time.Minute)

		for {
			select {
			case <-ticker.C:

				tool.Mutex.Lock()
				_, ok := tool.processes[name]
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

	return <-done

}

func executeTestProgram(sandbox string, program string, parameters []string) (string, error) {

	stat0, err0 := os.Stat(sandbox)

	if err0 == nil {

		if stat0.IsDir() == true {

			stat1, err1 := os.Stat(program)

			if err1 == nil {

				if stat1.IsDir() == false {

					arguments := append([]string{
						"run",
						program,
					}, parameters...)

					cmd := exec.Command("go", arguments...)
					cmd.Dir = sandbox

					output, err1 := cmd.Output()

					if err1 == nil {
						return strings.TrimSpace(string(output)), nil
					} else {
						return "", err1
					}

				} else {
					return "", fmt.Errorf("%s is not a file!", sandbox)
				}

			} else {
				return "", err1
			}

		} else {
			return "", fmt.Errorf("%s is not a directory!", sandbox)
		}

	} else {
		return "", err0
	}

}

func TestAgents_Hire(t *testing.T) {

	playground, _ := os.MkdirTemp("/tmp", "exocomp-test-agents-*")
	sandbox       := filepath.Join(playground, "agents")
	model         := "qwen3-coder:30b"
	url,        _ := net_url.Parse("http://localhost:11434/v1")
	tool          := NewAgents(playground, sandbox, model, url, true)

	if tool != nil {

		result0, err0 := tool.Hire(
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

		if result0 != "agents.Hire: Agent \"Fibonacci Worker\" got hired." {
			t.Errorf("Expected agent \"%s\" toe get hired", "Fibonacci Worker")
		}

		if err0 == nil {

			finished0 := waitForAgent(tool, "Fibonacci Worker")

			if finished0 == true {

				result1, err1 := executeTestProgram(sandbox + "/fibonacci", sandbox + "/fibonacci/main.go", []string{"1"})
				result2, err2 := executeTestProgram(sandbox + "/fibonacci", sandbox + "/fibonacci/main.go", []string{"2"})
				result3, err3 := executeTestProgram(sandbox + "/fibonacci", sandbox + "/fibonacci/main.go", []string{"3"})
				result4, err4 := executeTestProgram(sandbox + "/fibonacci", sandbox + "/fibonacci/main.go", []string{"4"})
				result5, err5 := executeTestProgram(sandbox + "/fibonacci", sandbox + "/fibonacci/main.go", []string{"5"})
				result6, err6 := executeTestProgram(sandbox + "/fibonacci", sandbox + "/fibonacci/main.go", []string{"6"})

				if result1 != "0" {
					t.Errorf("Expected %s to be %s for step %d", result1, "0", 1)
				}

				if err1 != nil {
					t.Errorf("Expected %v to be nil", err1)
				}
				
				if result2 != "1" {
					t.Errorf("Expected %s to be %s for step %d", result2, "1", 2)
				}

				if err2 != nil {
					t.Errorf("Expected %v to be nil", err2)
				}

				if result3 != "1" {
					t.Errorf("Expected %s to be %s for step %d", result3, "1", 3)
				}

				if err3 != nil {
					t.Errorf("Expected %v to be nil", err3)
				}

				if result4 != "2" {
					t.Errorf("Expected %s to be %s for step %d", result4, "2", 4)
				}

				if err4 != nil {
					t.Errorf("Expected %v to be nil", err4)
				}

				if result5 != "3" {
					t.Errorf("Expected %s to be %s for step %d", result5, "3", 5)
				}

				if err5 != nil {
					t.Errorf("Expected %v to be nil", err5)
				}

				if result6 != "5" {
					t.Errorf("Expected %s to be %s for step %d", result6, "5", 6)
				}

				if err6 != nil {
					t.Errorf("Expected %v to be nil", err6)
				}

			} else {
				t.Errorf("Expected agent \"%s\" to call agents.Quit", "Fibonacci Worker")
			}

		} else {
			t.Errorf("Expected %v to be nil", err0)
		}

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

func TestAgents_List(t *testing.T) {
	fmt.Println("TODO: Test agents.List()")
}

