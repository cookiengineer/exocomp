package tty

import "exocomp/schemas"
import "exocomp/tools"
import "exocomp/types"
import "bufio"
import "fmt"
import "os"
import "os/signal"
import "strings"
import "syscall"
import "time"

type Client struct {
	Renderer *Renderer
	Session  *types.Session
	Role     string
}

func NewClient(agent *types.Agent, config *types.Config) *Client {

	var session *types.Session = nil

	recovery := types.NewRecovery(config.Playground)

	if recovery.HasBackup() {

		session = recovery.RestoreSession()

		if session == nil {
			session = types.NewSession(agent, config)
		}

	} else {
		session = types.NewSession(agent, config)
	}

	renderer := NewRenderer(session)

	if len(agent.AllowedTools) > 0 {

		tool_schemas, tools := tools.Toolset(
			config.Playground,
			config.Sandbox,
			config.Model,
			config.URL,
			config.Debug,
			agent.AllowedPrograms,
			agent.AllowedTools,
		)

		for name, tool := range tools {
			session.SetTool(name, tool, tool_schemas[name])
		}

	}

	tool := session.GetTool("agents.List")

	if tool != nil {

		agent_tool, ok := tool.(*tools.Agents)

		if ok == true {

			agents := recovery.RestoreAgents()

			if len(agents) > 0 {

				for _, agent := range agents {
					agent_tool.SetAgent(agent)
				}

			}

		}

	}

	if config.GetPrompt() != "" {
		session.Init()
	}

	return &Client{
		Renderer: renderer,
		Session:  session,
		Role:     "user",
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
		signals<-syscall.SIGINT
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

		role   := client.Role
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
