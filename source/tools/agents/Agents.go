package agents

import agents_roles "exocomp/agents"
import "exocomp/schemas"
import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"
import "context"
import "fmt"
import "io"
import net_url "net/url"
import "os/exec"
import "slices"
import "sort"
import "sync"
import "time"

type agentProcess struct {
	cmd         *exec.Cmd
	cancel      context.CancelFunc
	done        chan struct{}
	last_write  time.Time
	stdin       io.WriteCloser
	write_mutex *sync.Mutex
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
	OnSwitch    func(role string) (string, error)
	mutex       *sync.Mutex
	contents    map[string]*types.Agent
	processes   map[string]*agentProcess
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
		OnSwitch:    nil,
		mutex:       &sync.Mutex{},
		contents:    make(map[string]*types.Agent),
		processes:   make(map[string]*agentProcess),
	}

	// NOTE: readAgents() allowed only at bootup time
	readAgents(agents)

	return agents

}

func (tool *Agents) Name() string {
	return "agents"
}

func (tool *Agents) Call(method string, arguments map[string]interface{}) (string, error) {

	if tool.HasMethod(method) == true {

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

		} else if method == "Switch" {

			role, ok1 := arguments["role"].(string)

			if ok1 == true {
				return tool.Switch(role)
			} else {
				return "", fmt.Errorf("agents.%s: %s", method, "Invalid parameter \"role\" is not a string.")
			}

		} else if method == "Hire" {

			role, ok1 := arguments["role"].(string)
			prompt, ok2 := arguments["prompt"].(string)
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

func (tool *Agents) GetContent(id string) (any, error) {

	name := utils_fmt.FormatAgentName(id)

	tool.mutex.Lock()
	content, ok := tool.contents[name]
	tool.mutex.Unlock()

	if ok == true {
		return content, nil
	} else {
		return nil, fmt.Errorf("agents.Get: %s does not exist?", name)
	}

}

func (tool *Agents) GetContentIdentifiers() []string {

	result := make([]string, 0)

	tool.mutex.Lock()

	for id, _ := range tool.contents {
		result = append(result, id)
	}

	tool.mutex.Unlock()

	sort.Strings(result)

	return result

}

func (tool *Agents) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Agents) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range AgentsSchema {

			// NOTE: Patch the JSON Schema based on available Agent Roles _at runtime_
			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), "Hire") || schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), "Switch") {

				roles := make([]string, 0)

				for role, _ := range agents_roles.Roles {

					if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), "Hire") && role == "planner" {
						continue
					}

					roles = append(roles, role)

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
