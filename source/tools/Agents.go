package tools

import "exocomp/agents"
import "exocomp/schemas"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import utils_fmt "exocomp/utils/fmt"
import utils_rand "exocomp/utils/rand"
import "bufio"
import "bytes"
import "context"
import "encoding/json"
import "fmt"
import "io"
import net_url "net/url"
import "os"
import "os/exec"
import "sort"
import "strings"
import "sync"
import "time"

type Agents struct {
	Playground string
	Sandbox    string
	Model      string
	URL        *net_url.URL
	Debug      bool
	Mutex      *sync.Mutex
	contents   map[string]*types.Agent
	processes  map[string]*os.Process
}

func NewAgents(playground string, sandbox string, model string, url *net_url.URL, debug bool) *Agents {

	agents := &Agents{
		Playground: playground,
		Sandbox:    sandbox,
		Model:      model,
		URL:        url,
		Debug:      debug,
		Mutex:      &sync.Mutex{},
		contents:   make(map[string]*types.Agent),
		processes:  make(map[string]*os.Process),
	}

	// NOTE: readAgents() allowed only at bootup time
	readAgents(agents)

	return agents

}

func (tool *Agents) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "List" {

		return tool.List()

	} else if method == "Roles" {

		return tool.Roles()

	} else if method == "Hire" {

		role,    ok1 := arguments["role"].(string)
		prompt,  ok2 := arguments["prompt"].(string)

		name,    ok3 := arguments["name"].(string)
		sandbox, ok4 := arguments["sandbox"].(string)

		if ok1 == true && ok2 == true && ok3 == true && ok4 == true {

			return tool.Hire(
				utils_fmt.FormatAgentRole(role),
				utils_fmt.FormatMultiLine(prompt),
				utils_fmt.FormatAgentName(name),
				sandbox,
			)

		} else if ok1 == true && ok2 == true && ok3 == true && ok4 == false {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"sandbox\" is not a string.")
		} else if ok1 == true && ok2 == true && ok3 == false && ok4 == true {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"name\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true && ok4 == true {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"prompt\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true && ok4 == true {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"role\" is not a string.")
		} else {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameters.")
		}

	} else if method == "Fire" {

		name, ok1 := arguments["name"].(string)

		if ok1 == true {
			return tool.Fire(utils_fmt.FormatAgentName(name))
		} else {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"name\" is not a string.")
		}

	} else if method == "Inquire" {

		name, ok1 := arguments["name"].(string)

		if ok1 == true {
			return tool.Inquire(utils_fmt.FormatAgentName(name))
		} else {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"name\" is not a string.")
		}

	} else if method == "Quit" {

		message, ok1 := arguments["message"].(string)

		if ok1 == true {
			return tool.Quit(utils_fmt.FormatMultiLine(message))
		} else {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"message\" is not a string.")
		}

	} else {
		return "", fmt.Errorf("agents.%s: Invalid method.", method)
	}

}

func (tool *Agents) Get(id string) (any, error) {

	name := utils_fmt.FormatAgentName(id)

	tool.Mutex.Lock()
	content, ok := tool.contents[name]
	tool.Mutex.Unlock()

	if ok == true {
		return content, nil
	}

	return nil, fmt.Errorf("agents.Get: %s does not exist?", name)

}

func (tool *Agents) GetAgent(id string) *types.Agent {

	name := utils_fmt.FormatAgentName(id)

	tool.Mutex.Lock()
	agent, ok := tool.contents[name]
	tool.Mutex.Unlock()

	if ok == true {
		return agent
	} else {
		return nil
	}

}

func (tool *Agents) GetNames() []string {

	result := make([]string, 0)

	tool.Mutex.Lock()
	for name, _ := range tool.contents {
		result = append(result, name)
	}
	tool.Mutex.Unlock()

	sort.Strings(result)

	return result

}

func (tool *Agents) List() (string, error) {

	len_content := 0

	tool.Mutex.Lock()
	len_content = len(tool.contents)
	tool.Mutex.Unlock()


	if len_content > 0 {

		lines := make([]string, 0)

		tool.Mutex.Lock()
		for name, agent := range tool.contents {

			_, ok  := tool.processes[name]
			status := "unknown"

			if ok == true {
				status = "working"
			} else {
				status = "finished"
			}

			lines = append(lines, fmt.Sprintf("- Name: \"%s\", Type: %s, Status: %s", agent.Name, agent.Role, status))

		}
		tool.Mutex.Unlock()

		sort.Strings(lines)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("agents.List: %d agents were working for us.", len(lines)))

		for l := 0; l < len(lines); l++ {
			result = append(result, lines[l])
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("agents.List: No agents are working for us!")
	}

}

func (tool *Agents) Roles() (string, error) {

	lines := make([]string, 0)

	for _, template := range agents.Roles {
		lines = append(lines, fmt.Sprintf("- Role: \"%s\", Description: %s", template.Role, template.Description))
	}

	sort.Strings(lines)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("agents.Roles: %d available agent roles.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}

func (tool *Agents) Hire(role string, prompt string, name string, sandbox string) (string, error) {

	if name == "" || name == "." {
		name = utils_rand.AgentName(role)
	}

	if sandbox == "" || sandbox == "." || sandbox == "./" {
		sandbox = tool.Sandbox
	}

	tool.Mutex.Lock()
	_, ok := tool.contents[name]
	tool.Mutex.Unlock()

	if ok == false {

		resolved, err0 := resolveSandboxPath(tool.Sandbox, sandbox)

		if err0 == nil {

			stat, err1 := os.Stat(resolved)

			if err1 == nil && stat.IsDir() == true {
				// Do Nothing
			} else {
				os.MkdirAll(resolved, 0755)
			}

			debug_flag := ""

			if tool.Debug == true {
				debug_flag = "--debug"
			}

			exe, _ := os.Executable()

			if os.Getenv("EXOCOMP_AGENT") != "" {
				exe = os.Getenv("EXOCOMP_AGENT")
			}

			// NOTE: child's playground is parent's sandbox
			ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Minute)
			cmd := exec.CommandContext(
				ctx,
				exe,
				"agent",
				fmt.Sprintf("--name=\"%s\"", name),
				fmt.Sprintf("--role=\"%s\"", role),
				fmt.Sprintf("--model=\"%s\"", tool.Model),
				fmt.Sprintf("--prompt=\"%s\"", prompt),
				// --temperature set by agent role
				fmt.Sprintf("--playground=\"%s\"", tool.Sandbox),
				fmt.Sprintf("--sandbox=\"%s\"", resolved),
				fmt.Sprintf("--url=\"%s\"", tool.URL.String()),
				debug_flag,
			)
			cmd.Dir = resolved

			cmd.Stdin = strings.NewReader("")

			// XXX: Use this for debugging
			// cmd.Stderr = os.Stderr

			stdout_pipe, err2 := cmd.StdoutPipe()

			if err2 == nil {

				err3 := cmd.Start()

				if err3 == nil {

					tool.Mutex.Lock()
					tool.contents[name] = agents.NewAgent(types.NewConfig(
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

					if debug_flag != "" {
						// XXX: "exocomp agent" prints first system message
						tool.contents[name].Messages = make([]*schemas.Message, 0)
					}

					tool.processes[name] = cmd.Process
					tool.Mutex.Unlock()



					// Background Reaper
					go func(name string, tool *Agents, stdout_pipe io.ReadCloser) {

						scanner := bufio.NewScanner(stdout_pipe)

						for scanner.Scan() {

							line := scanner.Text()

							if strings.HasPrefix(line, "schemas.Message:") {

								buffer  := []byte(line[16:])
								message := schemas.Message{}

								err3 := json.Unmarshal(buffer, &message)

								if err3 == nil {

									tool.Mutex.Lock()

									agent, ok1 := tool.contents[name]

									if ok1 == true {
										agent.Messages = append(agent.Messages, &message)
									}

									tool.Mutex.Unlock()

								}

							} else if strings.HasPrefix(line, "types.ContextUsage:") {

								buffer   := []byte(line[19:])
								usage    := types.ContextUsage{}

								err3 := json.Unmarshal(buffer, &usage)

								if err3 == nil {

									tool.Mutex.Lock()

									agent, ok1 := tool.contents[name]

									if ok1 == true {
										agent.ContextUsage.Length = usage.Length
										agent.ContextUsage.Tokens = usage.Tokens
									}

									tool.Mutex.Unlock()

								}

							}

						}

						stdout_pipe.Close()

					}(name, tool, stdout_pipe)

					// Background Reaper
					go func(name string, tool *Agents, ctx context.Context, cancel context.CancelFunc) {

						tool.Mutex.Lock()
						last_length := len(tool.contents[name].Messages)
						tool.Mutex.Unlock()

						last_time := time.Now()
						ticker    := time.NewTicker(10 * time.Second)

						BackgroundLoop:
						for {

							select {
							case <-ctx.Done():

								break BackgroundLoop

							case <-ticker.C:

								if time.Since(last_time) > 1 * time.Minute {

									cancel()
									break BackgroundLoop

								} else {

									tool.Mutex.Lock()
									length := len(tool.contents[name].Messages)
									tool.Mutex.Unlock()

									if length > last_length {
										last_length = length
										last_time   = time.Now()
									}

								}

							}

						}

						ticker.Stop()

					}(name, tool, ctx, cancel)

					// Background Reaper
					go func(name string, tool *Agents, cmd *exec.Cmd) {

						cmd.Wait()

						tool.Mutex.Lock()
						delete(tool.processes, name)
						tool.Mutex.Unlock()

					}(name, tool, cmd)

					sandbox_path, _ := sanitizeSandboxPath(resolved)

					return fmt.Sprintf("agents.Hire: Agent \"%s\" hired to work on \"%s\".", name, sandbox_path), nil

				} else {
					return "", fmt.Errorf("agents.Hire: %s", err3.Error())
				}

			} else {
				return "", fmt.Errorf("agents.Hire: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("agents.Hire: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("agents.Hire: Agent \"%s\" was already hired in the past. Pick a different name.", name)
	}

}

func (tool *Agents) Fire(name string) (string, error) {

	tool.Mutex.Lock()
	process, ok := tool.processes[name]
	tool.Mutex.Unlock()

	if ok == true {

		err := process.Kill()

		if err == nil {

			tool.Mutex.Lock()
			delete(tool.processes, name)
			tool.Mutex.Unlock()

			return fmt.Sprintf("agents.Fire: Agent \"%s\" fired.", name), nil

		} else {
			return "", fmt.Errorf("agents.Fire: %s", err.Error())
		}

	} else {
		return "", fmt.Errorf("agents.Fire: Agent \"%s\" already quit!", name)
	}

}

func (tool *Agents) Inquire(name string) (string, error) {

	tmp, err0 := os.MkdirTemp("/tmp", "exocomp-summarizer-*")

	if err0 == nil {

		tool.Mutex.Lock()
		agent, ok0 := tool.contents[name]
		tool.Mutex.Unlock()

		if ok0 == true {

			tool.Mutex.Lock()
			messages := utils_chat.SummarizeMessages(agent.Messages, true, true, false)
			tool.Mutex.Unlock()

			prompt := strings.Join([]string{
				"Please summarize the following conversation, the latest messages are the newest ones.",
				"",
				messages,
			}, "\n")

			debug_flag := ""

			if tool.Debug == true {
				debug_flag = "--debug"
			}

			exe, _ := os.Executable()

			if os.Getenv("EXOCOMP_AGENT") != "" {
				exe = os.Getenv("EXOCOMP_AGENT")
			}

			cmd := exec.Command(
				exe,
				"agent",
				fmt.Sprintf("--name=\"%s\"", "Summarizer"),
				fmt.Sprintf("--role=\"%s\"", "summarizer"),
				fmt.Sprintf("--model=\"%s\"", tool.Model),
				fmt.Sprintf("--prompt=\"%s\"", prompt),
				// --temperature set by agent role
				// --playground set by cmd.Dir
				// --sandbox set by cmd.Dir
				fmt.Sprintf("--url=\"%s\"", tool.URL.String()),
				debug_flag,
			)
			cmd.Dir = tmp

			stdout_buffer := bytes.Buffer{}
			cmd.Stdout = &stdout_buffer

			err1 := cmd.Run()

			if err1 == nil {

				os.RemoveAll(tmp)

				lines := strings.Split(strings.TrimSpace(stdout_buffer.String()), "\n")

				if len(lines) > 0 {

					summary := schemas.Message{}
					err2    := json.Unmarshal([]byte(lines[len(lines) - 1]), &summary)

					if err2 == nil {

						tool.Mutex.Lock()
						_, ok1 := tool.processes[name]
						tool.Mutex.Unlock()

						if ok1 == true {

							result := strings.Join([]string{
								fmt.Sprintf("agents.Inquire: Summary of currently working agent \"%s\"'s work report:", name),
								strings.TrimSpace(summary.Content),
							}, "\n")

							return result, nil

						} else {

							result := strings.Join([]string{
								fmt.Sprintf("agents.Inquire: Summary of already finished agent \"%s\"'s work report:", name),
								strings.TrimSpace(summary.Content),
							}, "\n")

							return result, nil

						}

					} else {
						return "", fmt.Errorf("agents.Inquire: Failed to summarize agent \"%s\"'s work report!", name)
					}

				} else {
					return "", fmt.Errorf("agents.Inquire: Failed to summarize agent \"%s\"'s work report!", name)
				}

			} else {
				return "", fmt.Errorf("agents.Inquire: Failed to summarize agent \"%s\"'s work report!", name)
			}

		} else {
			return "", fmt.Errorf("agents.Inquire: Agent \"%s\" didn't work for us!", name)
		}

	} else {
		return "", fmt.Errorf("agents.Inquire: System is out of memory ... %s", err0.Error())
	}

}

func (tool *Agents) Quit(message string) (string, error) {

	if strings.Contains(strings.ToLower(message), "my work is done") {

		go func() {
			// Give Renderer time to catch up
			time.Sleep(200 * time.Millisecond)
			os.Exit(0)
		}()

		return fmt.Sprintf("agents.Quit: Agent quit with work report:\n%s", strings.TrimSpace(message)), nil

	} else {

		go func() {
			// Give Renderer time to catch up
			time.Sleep(200 * time.Millisecond)
			os.Exit(1)
		}()

		return fmt.Sprintf("agents.Quit: Agent quit with work report:\n%s", strings.TrimSpace(message)), nil

	}

}

