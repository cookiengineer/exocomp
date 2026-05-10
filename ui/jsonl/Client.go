package jsonl

import "exocomp/schemas"
import "exocomp/tools"
import "exocomp/types"
import utils_cli "exocomp/utils/cli"
import "bufio"
import "encoding/json"
import "fmt"
import "os"
import "os/signal"
import "strings"
import "syscall"

type Client struct {
	Renderer *Renderer
	Session  *types.Session
	role     string
}

func NewClient(agent *types.Agent, config *types.Config) *Client {

	session  := types.NewSession(agent, config)
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

	if config.GetPrompt() != "" {
		session.Init()
	}

	return &Client{
		Renderer: renderer,
		Session:  session,
		role:     "user",
	}

}

func (client *Client) Init() {

	signals := make(chan os.Signal, 1)

	signal.Notify(
		signals,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		client.Session.Init()
	}()

	go func() {
		client.InputLoop()
		signals<-syscall.SIGINT
	}()

	go client.Renderer.RenderLoop()

	select {
	case sig := <-signals:

		switch sig {
		case syscall.SIGINT:

			client.Destroy()
			fmt.Fprintf(os.Stdout, "Received signal: %s\n", "SIGINT")
			os.Exit(0)

		case syscall.SIGTERM:

			client.Destroy()
			fmt.Fprintf(os.Stdout, "Received signal: %s\n", "SIGTERM")
			os.Exit(0)

		default:

			client.Destroy()
			fmt.Printf("Received signal: %s\n", sig.String())
			os.Exit(0)

		}

	}

}

func (client *Client) InputLoop() {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		role   := client.role
		prompt := strings.TrimSpace(scanner.Text())

		if prompt != "" && client.Session != nil {

			if role == "user" || role == "assistant" {

				if strings.HasPrefix(prompt, "/") && strings.Contains(prompt, " ") && !strings.Contains(prompt, "\n") {

					identifier := prompt[1:strings.Index(prompt, " ")]

					if strings.Contains(identifier, ".") {

						method    := identifier[strings.LastIndex(identifier, ".")+1:]
						arguments := utils_cli.ParseParameters(strings.TrimSpace(prompt[1+len(identifier)+1:]))

						client.Session.CallTool(identifier, method, arguments)

					}

				} else if strings.HasPrefix(prompt, "{") && strings.HasSuffix(prompt, "}") {

					tmp  := schemas.Message{}
					err1 := json.Unmarshal([]byte(prompt), &tmp)

					if err1 == nil && role == tmp.Role {

						go func() {

							err2 := client.Session.SendChatRequest(schemas.Message{
								Role:    tmp.Role,
								Content: tmp.Content,
							})

							if err2 != nil {
								os.Exit(1)
							}

						}()

					} else {
						fmt.Fprintf(os.Stderr, "Error: jsonl.Client: %s", "Invalid schemas.Message")
					}

				}

			}

		}

	}

}

func (client *Client) Destroy() {

	if client.Renderer != nil {
		client.Renderer.Destroy()
	}

}

func (client *Client) SetRole(role string) {

	if role == "user" || role == "assistant" {
		client.role = role
	}

}
