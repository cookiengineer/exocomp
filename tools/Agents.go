package tools

import "exocomp/agents"
import "exocomp/schemas"
import utils_fmt "exocomp/utils/fmt"
import "bufio"
import "encoding/json"
import "fmt"
import "os"
import "os/exec"
import "sort"
import "strings"
import "sync"

type Agents struct {
	Sandbox    string
	Playground string
	mutex      *sync.Mutex
	agents     map[string]*agents.Agent
	chats      map[string][]*schemas.Message
	processes  map[string]*os.Process
}

func NewAgents(agent string, sandbox string, playground string) *Agents {

	agents := &Agents{
		Sandbox:    sandbox,
		Playground: playground,
		mutex:      &sync.Mutex{},
		agents:     make(map[string]*agents.Agent),
		chats:      make(map[string][]*schemas.Message, 0),
		processes:  make(map[string]*os.Process),
	}

	return agents

}

func (tool *Agents) Call(method string, arguments map[string]interface{}) (string, error) {

	if method == "List" {

		return tool.List()

	} else if method == "Hire" {

		name,    ok1 := arguments["name"].(string)
		agent,   ok2 := arguments["agent"].(string)
		sandbox, ok3 := arguments["sandbox"].(string)
		prompt,  ok4 := arguments["prompt"].(string)

		if ok1 == true && ok2 == true && ok3 == true && ok4 == true {
			return tool.Hire(utils_fmt.FormatAgentName(name), agent, sandbox, utils_fmt.FormatMultiLine(prompt))
		} else if ok1 == true && ok2 == true && ok3 == true && ok4 == false {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"prompt\" is not a string.")
		} else if ok1 == true && ok2 == true && ok3 == false && ok4 == true {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"sandbox\" is not a string.")
		} else if ok1 == true && ok2 == false && ok3 == true && ok4 == true {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"agent\" is not a string.")
		} else if ok1 == false && ok2 == true && ok3 == true && ok4 == true {
			return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"name\" is not a string.")
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

func (tool *Agents) List() (string, error) {

	if len(tool.agents) > 0 {

		lines := make([]string, 0)

		for _, agent := range tool.agents {
			lines = append(lines, fmt.Sprintf("- \"%s\" (%s)", agent.Name, agent.Type.String()))
		}

		sort.Strings(lines)

		result := make([]string, 0)
		result = append(result, fmt.Sprintf("agents.List: %d agents working for us.", len(lines)))

		for l := 0; l < len(lines); l++ {
			result = append(result, lines[l])
		}

		return strings.Join(result, "\n"), nil

	} else {
		return "", fmt.Errorf("agents.List: No agents are working for us!")
	}

}

func (tool *Agents) Hire(name string, agent string, sandbox string, prompt string) (string, error) {

	resolved, err0 := resolveSandboxPath(tool.Sandbox, sandbox)

	if err0 == nil {

		cmd := exec.Command(
			os.Args[0],
			"jsonl",
			"--name",       name,
			"--agent",      agent,
			"--playground", tool.Playground,
			"--sandbox",    resolved,
			"--prompt",     prompt,
		)
		cmd.Dir = resolved


		stdout_pipe, err1 := cmd.StdoutPipe()

		if err1 == nil {

			err2 := cmd.Start()

			if err2 == nil {

				tool.agents[name]    = agents.NewAgent(name, agent, "", 0.0)
				tool.processes[name] = cmd.Process

				// Background Reader
				go func(name string) {

					scanner := bufio.NewScanner(stdout_pipe)

					for scanner.Scan() {

						line    := scanner.Bytes()
						message := schemas.Message{}

						err := json.Unmarshal(line, &message)

						if err == nil {

							tool.mutex.Lock()

							_, ok := tool.chats[name]

							if ok == false {
								tool.chats[name] = make([]*schemas.Message, 0)
							}

							tool.chats[name] = append(tool.chats[name], &message)

							tool.mutex.Unlock()

						}

					}

				}(name)

				// Background Reaper
				go func(name string, cmd *exec.Cmd) {

					cmd.Wait()

					tool.mutex.Lock()
					delete(tool.agents, name)
					delete(tool.processes, name)
					tool.mutex.Unlock()

				}(name, cmd)

				return fmt.Sprintf("agents.Hire: Agent \"%s\" got hired.", name), nil

			} else {
				return "", fmt.Errorf("agents.Hire: %s", err2.Error())
			}

		} else {
			return "", fmt.Errorf("agents.Hire: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("agents.Hire: %s", err0.Error())
	}

}

func (tool *Agents) Fire(name string) (string, error) {

	process, ok := tool.processes[name]

	if ok == true {

		err := process.Kill()

		if err == nil {

			tool.mutex.Lock()
			delete(tool.agents, name)
			delete(tool.processes, name)
			tool.mutex.Unlock()

			return fmt.Sprintf("agents.Fire: Agent \"%s\" got fired.", name), nil

		} else {
			return "", fmt.Errorf("agents.Fire: %s", err.Error())
		}

	} else {
		return "", fmt.Errorf("agents.Fire: Agent \"%s\" doesn't work for us anymore?", name)
	}

}

func (tool *Agents) Quit(message string) (string, error) {

	if strings.Contains(strings.ToLower(message), "my work is done") {
		os.Exit(0)
		return "Quitting...", nil
	} else {
		os.Exit(1)
		return "Quitting...", nil
	}

}

