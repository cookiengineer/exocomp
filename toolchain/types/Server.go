package types

import "bufio"
import "fmt"
import "net/url"
import "os"
import "os/exec"
import "strconv"
import "strings"

type Pipe interface {
	Read([]byte) (int, error)
}

type Server struct {
	URL     *url.URL
	Ready   chan bool
	Errors  chan error
	cmd     *exec.Cmd
	stderr  Pipe
	watched bool
}

func NewServer(listener_url *url.URL, model_name string, model_path string) *Server {

	cmd := exec.Command(
		"llama-server",
		"--model",       model_path,
		"--alias",       model_name,
		"--gpu-layers",  "all",
		"--ctx-size",    strconv.Itoa(32768),
		"--batch-size",  strconv.Itoa(512),
		"--ubatch-size", strconv.Itoa(128),
		"--flash-attn",  "auto",
		"--no-slots",
		"--no-webui",
		"--no-webui-mcp-proxy",
		"--jinja",
		"--host", listener_url.Hostname(),
		"--port", listener_url.Port(),
	)
	cmd.Dir = os.TempDir()

	stderr, err2 := cmd.StderrPipe()

	if err2 == nil {

		server := &Server{
			URL:     listener_url,
			cmd:     cmd,
			stderr:  stderr,
			Ready:   make(chan bool, 1),
			Errors:  make(chan error, 1),
			watched: false,
		}

		return server

	} else {
		return nil
	}

}

func (server *Server) Start() error {

	if server.watched == false {
		server.watched = true
		go server.watch()
	}

	if server.cmd.Process == nil {
		return server.cmd.Start()
	}

	return nil

}

func (server *Server) Stop() error {

	if server.cmd.Process != nil {
		return server.cmd.Process.Kill()
	}

	return nil

}

func (server *Server) watch() {

	scanner := bufio.NewScanner(server.stderr)

	for scanner.Scan() {

		line := scanner.Text()

		if strings.Contains(line, "vk::DeviceLostError") {

			server.Errors <- fmt.Errorf("%s", "vk::DeviceLostError")
			return

		} else if strings.Contains(line, "terminate called") {

			server.Errors <- fmt.Errorf("%s", line)
			return

		} else if strings.Contains(line, "main: server is listening") {

			server.Ready <- true
			return

		}

	}

	err := scanner.Err()

	if err != nil {
		server.Errors <- err
	}

}
