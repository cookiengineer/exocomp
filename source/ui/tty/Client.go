package tty

import "exocomp/adapters"
import "exocomp/agents"
import "exocomp/engine"
import "exocomp/schemas"
import "exocomp/tools"
import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import "fmt"

type Client struct {
	App     *App
	Session *engine.Session
	Role    string
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

	if agent.HasTools() {

		toolset := tools.Toolset(
			config.Playground,
			config.Sandbox,
			config.Model,
			config.URL,
			config.Debug,
			agent.AllowedPrograms,
			agent.AllowedImages,
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
		App:     NewApp(session),
		Session: session,
		Role:    "user",
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
								new_agent.AllowedImages,
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

							if client.App != nil {

								client.App.Transcript.Reset()
								client.App.refresh()

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

	if client.App != nil {
		client.App.Destroy()
	}

}

func (client *Client) Init() {

	if client.App != nil {

		defer client.Destroy()
		client.App.Init()

	}

}

func (client *Client) SetRole(role string) {

	if role == "user" || role == "assistant" {

		client.Role = role

		if client.App != nil {
			client.App.SetRole(role)
		}

	}

}
