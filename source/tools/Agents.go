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
import "slices"
import "sort"
import "strings"
import "sync"
import "time"

type agent_process struct {
	cmd        *exec.Cmd
	cancel     context.CancelFunc
	done       chan struct{}
	last_write time.Time
}

func handle_process_json_line(tool *Agents, name string, line []byte) {

	tool.Mutex.Lock()
	defer tool.Mutex.Unlock()

	process, ok := tool.processes[name]

	if ok == true {
		process.last_write = time.Now()
	}

	text := strings.TrimRight(string(line), "\n")

	if strings.HasPrefix(text, "schemas.Message:") {

		buffer  := []byte(text[16:])
		message := schemas.Message{}

		err := json.Unmarshal(buffer, &message)

		if err == nil {

			agent, ok1 := tool.contents[name]

			if ok1 == true {
				agent.Messages = append(agent.Messages, &message)
			}

		}

	} else if strings.HasPrefix(text, "types.ContextUsage:") {

		buffer := []byte(text[19:])
		usage  := types.ContextUsage{}

		err := json.Unmarshal(buffer, &usage)

		if err == nil {

			agent, ok1 := tool.contents[name]

			if ok1 == true {
				agent.ContextUsage.Length = usage.Length
				agent.ContextUsage.Tokens = usage.Tokens
			}

		}

	}

}

func read_process(tool *Agents, name string, process *agent_process, stdout_pipe io.ReadCloser) {

	reader := bufio.NewReader(stdout_pipe)

	for {

		line, err := reader.ReadBytes('\n')

		if len(line) > 0 {
			handle_process_json_line(tool, name, line)
		}

		if err != nil {
			break
		}

	}

	stdout_pipe.Close()

	process.cmd.Wait()

	tool.Mutex.Lock()

	agent, ok := tool.contents[name]

	if ok == true && agent.Status == "working" {

		if getAgentWorkReport(agent) != "" {
			agent.Status = "finished"
		} else {
			agent.Status = "failed"
		}

		agent.FinishedAt = schemas.NewDatetime()

	}

	delete(tool.processes, name)

	tool.Mutex.Unlock()

	close(process.done)

}

func watch_process(tool *Agents, process *agent_process) {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {

		select {

		case <-process.done:
			return

		case <-ticker.C:

			tool.Mutex.Lock()
			stale := tool.IdleTimeout > 0 && time.Since(process.last_write) > tool.IdleTimeout
			tool.Mutex.Unlock()

			if stale == true {
				process.cancel()
				return
			}

		}

	}

}

type Agents struct {
	Methods     []string
	Playground  string
	Sandbox     string
	Model       string
	URL         *net_url.URL
	Debug       bool
	Timeout     time.Duration
	IdleTimeout time.Duration
	OnQuit      func(report string, success bool)
	Mutex       *sync.Mutex
	contents    map[string]*types.Agent
	processes   map[string]*agent_process
}

func NewAgents(methods []string, playground string, sandbox string, model string, url *net_url.URL, debug bool) *Agents {

	agents := &Agents{
		Methods:     methods,
		Playground:  playground,
		Sandbox:     sandbox,
		Model:       model,
		URL:         url,
		Debug:       debug,
		Timeout:     30 * time.Minute,
		IdleTimeout: 5 * time.Minute,
		OnQuit:      nil,
		Mutex:       &sync.Mutex{},
		contents:    make(map[string]*types.Agent),
		processes:   make(map[string]*agent_process),
	}

	// NOTE: readAgents() allowed only at bootup time
	readAgents(agents)

	return agents

}

func (tool *Agents) Name() string {
	return "agents"
}

func (tool *Agents) Call(method string, arguments map[string]interface{}) (string, error) {

	if slices.Contains(tool.Methods, method) == true {

		if method == "Await" {

			name, ok1 := arguments["name"].(string)

			if ok1 == true {
				return tool.Await(utils_fmt.FormatAgentName(name))
			} else {
				return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"name\" is not a string.")
			}

		} else if method == "List" {

			return tool.List()

		} else if method == "Roles" {

			return tool.Roles()

		} else if method == "Hire" {

			role,    ok1 := arguments["role"].(string)
			prompt,  ok2 := arguments["prompt"].(string)
			sandbox, ok3 := arguments["sandbox"].(string)

			if ok1 == true && role == "planner" {

				return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"role\" can not be \"planner\".")

			} else if ok1 == true && ok2 == true && ok3 == true {

				return tool.Hire(
					utils_fmt.FormatAgentRole(role),
					utils_fmt.FormatMultiLine(prompt),
					sandbox,
				)

			} else if ok1 == true && ok2 == true && ok3 == false {

				return tool.Hire(
					utils_fmt.FormatAgentRole(role),
					utils_fmt.FormatMultiLine(prompt),
					".",
				)

			} else if ok1 == true && ok2 == false && ok3 == true {
				return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"prompt\" is not a string.")
			} else if ok1 == false && ok2 == true && ok3 == true {
				return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"role\" is not a string.")
			} else {
				return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameters \"role\" is not a string and \"prompt\" is not a string.")
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

	} else {
		return "", fmt.Errorf("agents.%s: Method not allowed.", method)
	}

}

// NOTE: Await blocks until the hired agent finished its work. This is
// deliberate: the planner model would otherwise poll in a hot loop and
// blow up its limited context window with identical "still working" tool
// messages. The lifecycle guarantees the done channel always closes
// (reader EOF -> cmd.Wait() -> finalize, forced by the liveness watchdog).
func (tool *Agents) Await(name string) (string, error) {

	tool.Mutex.Lock()
	process, running := tool.processes[name]
	tool.Mutex.Unlock()

	if running == true {
		<-process.done
	}

	tool.Mutex.Lock()
	agent, ok := tool.contents[name]
	report := ""
	status := ""

	if ok == true {
		report = getAgentWorkReport(agent)
		status = agent.Status
	}

	tool.Mutex.Unlock()

	if ok == false {
		return "", fmt.Errorf("agents.Await: Agent \"%s\" didn't work for us!", name)
	}

	if report != "" {
		return fmt.Sprintf("agents.Await: Agent \"%s\"'s Work Report\n===%s\n===", name, report), nil
	}

	if status == "fired" {
		return "", fmt.Errorf("agents.Await: Agent \"%s\" was fired!", name)
	} else {
		return "", fmt.Errorf("agents.Await: Agent \"%s\" never finished with a work report!", name)
	}

}

func (tool *Agents) GetContent(id string) (any, error) {

	name := utils_fmt.FormatAgentName(id)

	tool.Mutex.Lock()
	content, ok := tool.contents[name]
	tool.Mutex.Unlock()

	if ok == true {
		return content, nil
	}

	return nil, fmt.Errorf("agents.Get: %s does not exist?", name)

}

func (tool *Agents) GetAgent(name string) *types.Agent {

	sanitized_name := utils_fmt.FormatAgentName(name)

	tool.Mutex.Lock()
	agent, ok := tool.contents[sanitized_name]
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

			status := agent.Status

			if status == "" {
				status = "unknown"
			}

			lines = append(lines, fmt.Sprintf("- Name: \"%s\", Type: %s, Status: %s", name, agent.Role, status))

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

		if template.Role != "planner" {
			lines = append(lines, fmt.Sprintf("- Role: \"%s\", Description: %s", template.Role, template.Description))
		}

	}

	sort.Strings(lines)

	result := make([]string, 0)
	result = append(result, fmt.Sprintf("agents.Roles: %d available agent roles.", len(lines)))

	for l := 0; l < len(lines); l++ {
		result = append(result, lines[l])
	}

	return strings.Join(result, "\n"), nil

}

func (tool *Agents) Hire(role string, prompt string, sandbox string) (string, error) {

	name := utils_rand.AgentName(role)

	tool.Mutex.Lock()
	_, ok := tool.contents[name]
	tool.Mutex.Unlock()

	if ok == true {
		return "", fmt.Errorf("agents.Hire: Agent \"%s\" was already hired in the past. Pick a different name.", name)
	}

	resolved, err0 := resolveSandboxPath(tool.Sandbox, sandbox)

	if err0 != nil {
		return "", fmt.Errorf("agents.Hire: %s", err0.Error())
	}

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

	timeout := tool.Timeout

	if timeout <= 0 {
		timeout = 30 * time.Minute
	}

	// NOTE: child's playground is parent's sandbox
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

	if err2 != nil {
		cancel()
		return "", fmt.Errorf("agents.Hire: %s", err2.Error())
	}

	err3 := cmd.Start()

	if err3 != nil {
		cancel()
		return "", fmt.Errorf("agents.Hire: %s", err3.Error())
	}

	agent := agents.NewAgent(types.NewConfig(
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
	agent.Status    = "working"
	agent.StartedAt = schemas.NewDatetime()

	if debug_flag != "" {
		// XXX: "exocomp agent" prints first system message
		agent.Messages = make([]*schemas.Message, 0)
	}

	process := &agent_process{
		cmd:        cmd,
		cancel:     cancel,
		done:       make(chan struct{}),
		last_write: time.Now(),
	}

	tool.Mutex.Lock()
	tool.contents[name]  = agent
	tool.processes[name] = process
	tool.Mutex.Unlock()

	// NOTE: Single reader goroutine owns the child stdout lifecycle. It
	// reads unbounded lines, appends messages, and on EOF reaps the
	// process and finalizes the agent status.
	go read_process(tool, name, process, stdout_pipe)

	// NOTE: Liveness watchdog kills a child that stops producing output.
	go watch_process(tool, process)

	sandbox_path, _ := sanitizeSandboxPath(tool.Sandbox, resolved)

	return fmt.Sprintf("agents.Hire: Agent \"%s\" hired to work on \"%s\".", name, sandbox_path), nil

}

func (tool *Agents) Fire(name string) (string, error) {

	tool.Mutex.Lock()
	process, running := tool.processes[name]
	agent, ok     := tool.contents[name]

	if ok == true && running == true {
		agent.Status = "fired"
	}

	tool.Mutex.Unlock()

	if running == true {

		process.cancel()
		<-process.done

		return fmt.Sprintf("agents.Fire: Agent \"%s\" fired.", name), nil

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

			timeout := tool.Timeout

			if timeout <= 0 {
				timeout = 30 * time.Minute
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			cmd := exec.CommandContext(
				ctx,
				exe,
				"agent",
				fmt.Sprintf("--name=\"%s\"", utils_rand.AgentName("summarizer")),
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

	report := fmt.Sprintf("agents.Quit: Work Report\n%s", strings.TrimSpace(message))

	if strings.Contains(strings.ToLower(message), "my work is done") {

		if tool.OnQuit != nil {
			go tool.OnQuit(strings.TrimSpace(message), true)
		}

		return report, nil

	} else {

		if tool.OnQuit != nil {
			go tool.OnQuit(strings.TrimSpace(message), false)
		}

		return report, nil

	}

}

func (tool *Agents) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range AgentsSchema {

			// NOTE: Patch the JSON Schema based on available Agent Roles _at runtime_
			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), "Hire") {

				roles := make([]string, 0)

				for role, _ := range agents.Roles {

					if role != "planner" {
						roles = append(roles, role)
					}

				}

				sort.Strings(roles)

				property, ok := schema.Function.Parameters.Properties["role"]

				if ok == true {
					property.Enum = roles
					schema.Function.Parameters.Properties["role"] = property
				}

			}

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}

func (tool *Agents) SetAgent(agent *types.Agent) bool {

	sanitized_name := utils_fmt.FormatAgentName(agent.Name)

	if sanitized_name == agent.Name {

		tool.Mutex.Lock()
		_, ok := tool.contents[sanitized_name]
		tool.Mutex.Unlock()

		if ok == false {

			tool.Mutex.Lock()
			tool.contents[sanitized_name] = agent
			tool.Mutex.Unlock()

			return true

		}

	}

	return false

}

