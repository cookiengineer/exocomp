package tty

import "exocomp/adapters"
import "exocomp/agents"
import "exocomp/engine"
import "exocomp/schemas"
import "exocomp/tools"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "bufio"
import "fmt"
import "os"
import "os/signal"
import "strings"
import "syscall"
import "time"

type Client struct {
	Renderer *Renderer
	Session  *engine.Session
	Role     string
}

func NewClient(agent *types.Agent, config *types.Config) *Client {

	var session *engine.Session = nil

	recovery := engine.NewRecovery(config.Playground)

	if recovery.HasBackup() {

		session = recovery.RestoreSession()

		if session != nil {
			session.Console.Info("Restored Session from Backup")
		} else {
			session = engine.NewSession(agent, config)
			session.Console.Warn("Could not restore Session from Backup")
		}

	} else {
		session = engine.NewSession(agent, config)
	}

	renderer := NewRenderer(session)

	if agent.HasTools() {

		toolset := tools.Toolset(
			config.Playground,
			config.Sandbox,
			config.Model,
			config.URL,
			config.Debug,
			agent.AllowedPrograms,
			agent.AllowedTools,
		)

		if len(toolset) > 0 {

			for _, tool := range toolset {
				session.SetTool(tool)
			}

		}

	}

	if config.HasProvider(agent.Model) {

		adapterset := adapters.Adapterset(
			config.ResolveURL(agent.Model, ""),
			config.ResolveModel(agent.Model),
		)

		if len(adapterset) > 0 {

			for _, adapter := range adapterset {
				session.SetAdapter(adapter)
			}

		}

	}

	tool := session.GetTool("agents.List")
	agents := recovery.RestoreAgents()

	if tool != nil && len(agents) > 0 {

		agents_tool, ok := tool.(*tools.Agents)

		if ok == true {

			for _, agent := range agents {
				agents_tool.SetAgent(agent)
			}

		}

	}

	if config.GetPrompt() != "" {
		session.Init()
	}

	client := &Client{
		Renderer: renderer,
		Session:  session,
		Role:     "user",
	}

	client.bindEvents()

	return client

}

func RestoreClient(session *engine.Session, agents []*types.Agent) *Client {

	renderer := NewRenderer(session)

	if session.Agent.HasTools() {

		toolset := tools.Toolset(
			session.Config.Playground,
			session.Config.Sandbox,
			session.Config.Model,
			session.Config.URL,
			session.Config.Debug,
			session.Agent.AllowedPrograms,
			session.Agent.AllowedTools,
		)

		if len(toolset) > 0 {

			for _, tool := range toolset {
				session.SetTool(tool)
			}

		}

	}

	if session.Config.HasProvider(session.Agent.Model) {

		adapterset := adapters.Adapterset(
			session.Config.ResolveURL(session.Agent.Model, ""),
			session.Config.ResolveModel(session.Agent.Model),
		)

		if len(adapterset) > 0 {

			for _, adapter := range adapterset {
				session.SetAdapter(adapter)
			}

		}

	}

	tool := session.GetTool("agents.List")

	if tool != nil && len(agents) > 0 {

		agents_tool, ok := tool.(*tools.Agents)

		if ok == true {

			for _, agent := range agents {
				agents_tool.SetAgent(agent)
			}

		}

	}

	client := &Client{
		Renderer: renderer,
		Session:  session,
		Role:     "user",
	}

	client.bindEvents()

	return client

}

func (client *Client) bindEvents() {

	if client.Session != nil {

		tool := client.Session.GetTool("agents.Switch")

		if tool != nil {

			agents_tool, ok := tool.(*tools.Agents)

			if ok == true {
				agents_tool.OnSwitch = client.OnSwitch
			}

		}

	}

}

func (client *Client) OnSwitch(role string) (string, error) {

	if client.Session != nil {

		tool := client.Session.GetTool("agents.Switch")

		if tool != nil {

			agents_tool, ok0 := tool.(*tools.Agents)

			if ok0 == true {

				messages := client.Session.GetMessages(0)
				history := make([]*schemas.Message, 0)

				for m := 0; m < len(messages); m++ {

					message := messages[m]

					if message != nil && message.Role == "assistant" {

						tool_calls, err := utils_chat.FindToolCalls(message, []string{"agents.Switch"})

						if err == nil && len(tool_calls) > 0 {
							break
						}

					}

					history = append(history, message)

				}

				summary, err0 := agents_tool.Summarize(history)

				if err0 == nil {

					config := *client.Session.Config

					config.Name = ""
					config.Role = role
					config.Model = ""
					config.Prompt = ""
					config.Temperature = 0.0
					config.Providers = make(map[string]types.Provider)

					for model, provider := range client.Session.Config.Providers {
						config.Providers[model] = provider
					}

					if types.GlobalConfig != nil {

						for model, provider := range types.GlobalConfig.Providers {

							_, ok := config.Providers[model]

							if ok == false {
								config.Providers[model] = provider
							}

						}

					}

					new_agent := agents.NewAgent(&config)

					if new_agent != nil {

						if client.Session.Config.Providers == nil {
							client.Session.Config.Providers = make(map[string]types.Provider)
						}

						for model, provider := range config.Providers {
							client.Session.Config.Providers[model] = provider
						}

						err1 := client.Session.SwitchAgent(new_agent, summary)

						if err1 == nil {

							toolset := tools.Toolset(
								client.Session.Config.Playground,
								client.Session.Config.Sandbox,
								new_agent.Model,
								client.Session.Config.URL,
								client.Session.Config.Debug,
								new_agent.AllowedPrograms,
								new_agent.AllowedTools,
							)

							client.Session.SetTools(toolset)

							adapterset := make([]types.Adapter, 0)

							if client.Session.Config.HasProvider(new_agent.Model) {

								adapterset = adapters.Adapterset(
									client.Session.Config.ResolveURL(new_agent.Model, ""),
									client.Session.Config.ResolveModel(new_agent.Model),
								)

							}

							client.Session.SetAdapters(adapterset)

							client.bindEvents()

							if client.Renderer != nil {
								client.Renderer.Reset()
							}

							return fmt.Sprintf("agents.Switch: Switched to agent role \"%s\".", new_agent.Role), nil

						} else {
							return "", fmt.Errorf("agents.Switch: %s", err1.Error())
						}

					} else {
						return "", fmt.Errorf("agents.Switch: Unknown agent role \"%s\".", role)
					}

				} else {
					return "", err0
				}

			} else {
				return "", fmt.Errorf("agents.Switch: Tool does not exist.")
			}

		} else {
			return "", fmt.Errorf("agents.Switch: Tool does not exist.")
		}

	} else {
		return "", fmt.Errorf("agents.Switch: No session available.")
	}

}

func (client *Client) Destroy() {

	if client.Session != nil {

		client.Session.Recovery.BackupSession(client.Session)

		tool := client.Session.GetTool("agents.List")

		if tool != nil {

			agent_tool, ok := tool.(*tools.Agents)

			if ok == true {

				agent_names := agent_tool.GetNames()

				for _, name := range agent_names {

					agent := agent_tool.GetAgent(name)

					if agent != nil {
						client.Session.Recovery.BackupAgent(agent)
					}

				}

			}

		}

	}

	if client.Renderer != nil {
		client.Renderer.Destroy()
	}

}

func (client *Client) Init() {

	signals := make(chan os.Signal, 1)

	signal.Notify(
		signals,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	if client.Session != nil {
		go client.Session.Init()
	}

	go func() {
		client.InputLoop()
		signals <- syscall.SIGINT
	}()

	if client.Renderer != nil {
		go client.Renderer.RenderLoop()
	}

	select {
	case sig := <-signals:

		switch sig {
		case syscall.SIGINT:

			client.Destroy()
			fmt.Fprintf(os.Stdout, "Received signal: %s\n", "SIGINT")

			time.Sleep(1 * time.Second)
			os.Exit(0)

		case syscall.SIGTERM:

			client.Destroy()
			fmt.Fprintf(os.Stdout, "Received signal: %s\n", "SIGTERM")

			time.Sleep(1 * time.Second)
			os.Exit(0)

		default:

			client.Destroy()
			fmt.Printf("Received signal: %s\n", sig.String())

			time.Sleep(1 * time.Second)
			os.Exit(0)

		}

	}

}

func (client *Client) InputLoop() {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		role := client.Role
		prompt := strings.TrimSpace(scanner.Text())

		if prompt != "" && client.Session != nil {

			if role == "user" || role == "assistant" {

				if strings.HasPrefix(prompt, "/") && strings.Contains(prompt, " ") && !strings.Contains(prompt, "\n") {

					command := types.ParseCommand(prompt)

					if command != nil {
						client.Session.CallTool("", command.Name, command.Method, command.Arguments)
					}

				} else {

					go func() {

						err := client.Session.SendChatRequest(schemas.Message{
							Role:    role,
							Content: prompt,
						})

						if err != nil {

							fmt.Fprintf(os.Stderr, "\nFatal Error: %s\n", err.Error())
							os.Exit(1)

						}

					}()

				}

			}

		}

	}

}

func (client *Client) SetRole(role string) {

	if role == "user" || role == "assistant" {
		client.Role = role
		client.Renderer.Role = role
	}

}
