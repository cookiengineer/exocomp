package tty

import "exocomp/engine"
import "exocomp/types"
import "fmt"

type Client struct {
	App     *App
	Session *engine.Session
	Role    string
}

func NewClient(agent *types.Agent, config *types.Config) *Client {

	session, _ := engine.NewSessionWithRecovery(agent, config)

	session.InstallTools()
	session.InstallAdapters()

	agents := session.Recovery.RestoreAgents()

	session.InstallAgents(agents)

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
		client.Session.BindAgents(client.OnSwitch)
	}

}

func (client *Client) OnSwitch(role string) (string, error) {

	if client.Session != nil {

		result, err := client.Session.SwitchRole(role)

		if err == nil {

			client.bindEvents()

			if client.App != nil {
				client.App.Transcript.Reset()
				client.App.refresh()
			}

		}

		return result, err

	} else {
		return "", fmt.Errorf("agents.Switch: No session available.")
	}

}

func (client *Client) Destroy() {

	if client.Session != nil {
		client.Session.Backup()
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
