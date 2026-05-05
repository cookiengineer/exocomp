//go:build agents

package tools

import "exocomp/agents"
import ui_jsonl "exocomp/ui/jsonl"
import utils_cli "exocomp/utils/cli"
import "bufio"
import "context"
import "encoding/json"
import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "runtime"
import "strings"
import "testing"
import "time"

type Pipe interface {
	Read([]byte) (int, error)
}

func appendDebugLog(path string, bytes []byte) error {

	fd, err0 := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err0 == nil {

		defer fd.Close()

		fd.Write([]byte("TestMain:\n"))
		fd.Write(bytes)
		fd.Write([]byte("\n\n\n"))

		return nil

	} else {
		return err0
	}

}


func getToolsPath() string {

	_, filename, _, ok := runtime.Caller(0)

	if ok == true {
		return filepath.Dir(filename)
	} else {
		panic("Cannot get current test file path")
	}

}

func parseArgumentsFromProc() ([]string) {

	result := make([]string, 0)

	pid        := os.Getpid()
	bytes, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))

	if err == nil {

		arguments := strings.Split(string(bytes), "\x00")

		if len(arguments) > 0 && arguments[len(arguments)-1] == "" {
			arguments = arguments[0:len(arguments)-1]
		}

		for _, argument := range arguments {
			result = append(result, argument)
		}

	}

	return result

}

func watchServerOutput(pipe Pipe, ready chan bool, errors chan error) {

	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {

		line := scanner.Text()

		if strings.Contains(line, "vk::DeviceLostError") {

			errors <- fmt.Errorf("%s", "vk::DeviceLostError")
			return

		} else if strings.Contains(line, "terminate called") {

			errors <- fmt.Errorf("%s", line)
			return

		} else if strings.Contains(line, "main: server is listening") {

			ready <- true
			return

		}

	}

	err := scanner.Err()

	if err != nil {
		errors <- err
	}

}

func TestMain(main *testing.M) {

	actual_os_args := parseArgumentsFromProc()
	tmp, _ := json.MarshalIndent(actual_os_args, "", "\t")
	appendDebugLog("/tmp/whatthefuck", tmp)

	if len(actual_os_args) > 1 && actual_os_args[1] == "jsonl" {

		// NOTE: TestMain needs to reimplement main() to be able to spawn subprocesses
		config := utils_cli.ParseConfig(actual_os_args[1:])
		agent  := agents.NewAgent(config)

		config.Update(
			agent.Name,
			agent.Type,
			agent.Model,
			agent.Prompt,
			agent.Temperature,
		)

		err0 := os.MkdirAll(config.Sandbox, 0755)

		if err0 == nil {

			os.Stdout.Sync()

			client := ui_jsonl.NewClient(agent, config)
			client.Init()

		} else {
			os.Exit(1)
		}

	} else {

		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)
		defer cancel()

		tools_path  := getToolsPath()
		server_path := filepath.Join(tools_path, "..", "third_party/llama/llama-server")
		model_path  := filepath.Join(tools_path, "..", "third_party/models/qwen3-coder-30b-a3b-instruct-q8_0.gguf")

		cmd := exec.Command(
			server_path,
			"-m", model_path,
			"--gpu-layers", "all",
			"--ctx-size", "32768",
			"--batch-size", "512",
			"--ubatch-size", "128",
			"--cache-type-k", "q8_0",
			"--cache-type-v", "q8_0",
			"--flash-attn", "auto",
			"--no-slots",
			"--no-webui",
			"--no-webui-mcp-proxy",
			"--jinja",
			"--port", "11434",
		)
		cmd.Dir = "/tmp"

		stderr, err1 := cmd.StderrPipe()

		fmt.Println("==> Start llama-server")

		err2 := cmd.Start()

		if err1 == nil && err2 == nil {

			ready  := make(chan bool, 1)
			errors := make(chan error, 1)

			fmt.Println("--- Wait for llama-server ...")
			go watchServerOutput(stderr, ready, errors)

			select {
			case <-ready:

				fmt.Println("--- llama-server is ready!")

				code := main.Run()

				if cmd != nil && cmd.Process != nil {
					cmd.Process.Kill()
					cmd.Process.Wait()
				}

				os.Exit(code)

			case err := <- errors:

				panic(fmt.Sprintf("Llama server error: %v", err))

			case <-ctx.Done():
				panic("Llama server timeout")
			}

		} else {

			fmt.Println("--- llama-server exited with error")

			if err1 != nil {
				panic(err1)
			} else if err2 != nil {
				panic(err2)
			}

		}

	}

}
