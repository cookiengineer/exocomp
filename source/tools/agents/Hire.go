package agents

import agents_roles "exocomp/agents"
import "exocomp/schemas"
import "exocomp/types"
import utils_cli "exocomp/utils/cli"
import utils_rand "exocomp/utils/rand"
import utils_sandbox "exocomp/utils/sandbox"
import "context"
import "fmt"
import "os"
import "os/exec"
import "sync"
import "time"

func (tool *Agents) Hire(role string, prompt string, sandbox string) (string, error) {

	name := utils_rand.AgentName(role)

	tool.mutex.Lock()
	_, ok := tool.contents[name]
	tool.mutex.Unlock()

	if ok == false {

		resolved, err0 := utils_sandbox.ResolvePath(tool.Sandbox, sandbox)

		if err0 == nil {

			stat, err1 := os.Stat(resolved)

			if err1 == nil && stat.IsDir() == true {
				// Do Nothing
			} else {
				os.MkdirAll(resolved, 0755)
			}

			exe, err2 := utils_cli.FindAgentBinary()

			if err2 == nil {

				arguments := []string{
					fmt.Sprintf("--name=\"%s\"", name),
					fmt.Sprintf("--role=\"%s\"", role),
					fmt.Sprintf("--prompt=\"%s\"", prompt),
					// --temperature set by agent role
					fmt.Sprintf("--playground=\"%s\"", tool.Sandbox),
					fmt.Sprintf("--sandbox=\"%s\"", resolved),
					fmt.Sprintf("--url=\"%s\"", tool.URL.String()),
				}

				if tool.Debug == true {
					arguments = append(arguments, "--debug")
				}

				timeout := tool.Timeout

				if timeout <= 0 {
					timeout = 30 * time.Minute
				}

				// NOTE: child's playground is parent's sandbox
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				cmd := exec.CommandContext(ctx, exe, arguments...)
				cmd.Dir = resolved

				stdin_pipe, err3 := cmd.StdinPipe()

				if err3 == nil {

					// XXX: Use this for debugging
					// cmd.Stderr = os.Stderr

					stdout_pipe, err4 := cmd.StdoutPipe()

					if err4 == nil {

						err5 := cmd.Start()

						if err5 == nil {

							agent := agents_roles.NewAgent(types.NewConfig(
								name,
								role,
								tool.Model,
								prompt,
								0.0, // Don't change temperature
								tool.Sandbox,
								resolved,
								tool.URL,
								false,
							))

							if tool.Debug == true {
								// NOTE: exocomp-agent prints the first system message
								agent.Messages = make([]*schemas.Message, 0)
							}

							process := &agentProcess{
								cmd:         cmd,
								cancel:      cancel,
								done:        make(chan struct{}),
								last_write:  time.Now(),
								stdin:       stdin_pipe,
								write_mutex: &sync.Mutex{},
							}

							tool.mutex.Lock()
							tool.contents[name] = agent
							tool.processes[name] = process
							tool.mutex.Unlock()

							// NOTE: Single reader goroutine owns the child stdout lifecycle. It
							// reads unbounded lines, appends messages, and on EOF reaps the
							// process and finalizes the agent status.
							go readProcess(tool, name, process, stdout_pipe)

							// NOTE: Liveness watchdog kills a child that stops producing output.
							go watchProcess(tool, process)

							sandbox_path, _ := utils_sandbox.SanitizePath(tool.Sandbox, resolved)

							return fmt.Sprintf("agents.Hire: Agent \"%s\" hired to work in \"%s\" on:\n%s", name, sandbox_path, prompt), nil

						} else {

							stdin_pipe.Close()
							cancel()
							return "", fmt.Errorf("agents.Hire: Cannot start agent \"%s\".", name)

						}

					} else {

						stdin_pipe.Close()
						cancel()
						return "", fmt.Errorf("agents.Hire: Cannot capture agent output for \"%s\".", name)

					}

				} else {

					cancel()
					return "", fmt.Errorf("agents.Hire: Cannot capture agent input for \"%s\".", name)

				}

			} else {
				return "", fmt.Errorf("agents.Hire: exocomp-agent is not installed.")
			}

		} else {
			return "", fmt.Errorf("agents.Hire: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("agents.Hire: Agent \"%s\" was already hired in the past. Pick a different name.", name)
	}

}
