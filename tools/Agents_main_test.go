//go:build with_agents

package tools

import "bufio"
import "context"
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

func getToolsPath() string {

	_, filename, _, ok := runtime.Caller(0)

	if ok == true {
		return filepath.Dir(filename)
	} else {
		panic("Cannot get current test file path")
	}

}
func watchServerOutput(pipe Pipe, ready chan bool, errors chan error) {

	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {

		line := scanner.Text()

		if strings.Contains(line, "terminate called") || strings.Contains(line, "vk::DeviceLostError") {

			errors <- fmt.Errorf(line)
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

	stdout, err1 := cmd.StdoutPipe()
	// stderr, err_stderr := cmd.StderrPipe()

	err2 := cmd.Start()

	if err1 == nil && err2 == nil {

		ready  := make(chan bool, 1)
		errors := make(chan error, 1)

		go watchServerOutput(stdout, ready, errors)

		select {
		case <-ready:

			fmt.Println("Llama server is ready...")

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

		if err1 != nil {
			panic(err1)
		} else if err2 != nil {
			panic(err2)
		}

	}

}
