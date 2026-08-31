package jsonl

import "exocomp/adapters"
import "exocomp/engine"
import "exocomp/schemas"
import "exocomp/tools"
import "exocomp/types"
import "bufio"
import "encoding/json"
import "fmt"
import "os"
import "os/signal"
import "strings"
import "syscall"
import "time"

type Client struct {
	Renderer *Renderer
	Session  *engine.Session
	role     string
}

func NewClient(agent *types.Agent, config *types.Config) *Client {

	// NOTE: jsonl.Client has no Recovery

	session  := engine.NewSession(agent, config)
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

	client := &Client{
		Renderer: renderer,
		Session:  session,
		role:     "user",
	}

	agent_tool := session.GetTool("agents.List")

	if agent_tool != nil {

		agents_tool, ok := agent_tool.(*tools.Agents)

		if ok == true {
			agents_tool.OnQuit = client.QuitHook
		}

	}

	return client

}

func (client *Client) QuitHook(report string, success bool) {

	deadline := time.Now().Add(5 * time.Second)

	for time.Now().Before(deadline) {

		if client.Renderer.Flushed() == true {
			break
		}

		time.Sleep(20 * time.Millisecond)

	}

	os.Stdout.Sync()

	if success == true {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

}

func (client *Client) Destroy() {

	if client.Session != nil {
		client.Session.Destroy()
	}

	if client.Renderer != nil {
		client.Renderer.Destroy()
	}

}

func (client *Client) Init() bool {

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
	}()

	go func() {
		client.ContextUsageLoop()
	}()

	go func() {
		client.Renderer.RenderLoop()
	}()

	select {
	case sig := <-signals:

		switch sig {
		case syscall.SIGINT:

			client.Destroy()
			fmt.Fprintf(os.Stderr, "Received signal: %s\n", "SIGINT")
			os.Exit(0)

		case syscall.SIGTERM:

			client.Destroy()
			fmt.Fprintf(os.Stderr, "Received signal: %s\n", "SIGTERM")
			os.Exit(0)

		default:

			client.Destroy()
			fmt.Fprintf(os.Stderr, "Received signal: %s\n", sig.String())
			os.Exit(0)

		}

	}

	return true

}

func (client *Client) InputLoop() {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		role   := client.role
		prompt := strings.TrimSpace(scanner.Text())

		if prompt != "" && client.Session != nil {

			if role == "user" || role == "assistant" {

				if strings.HasPrefix(prompt, "/") && strings.Contains(prompt, " ") && !strings.Contains(prompt, "\n") {

					command := types.ParseCommand(prompt)

					if command != nil {
						client.Session.CallTool("", command.Name, command.Method, command.Arguments)
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

func (client *Client) ContextUsageLoop() {

	last_tokens := 0
	heartbeat   := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	for {

		select {

		case <-heartbeat.C:

			client.RenderContextUsage()

		default:

			tokens := client.Session.Agent.ContextUsage.Tokens

			if last_tokens != tokens {
				last_tokens = tokens
				client.RenderContextUsage()
			}

			time.Sleep(100 * time.Millisecond)

		}

	}

}

func (client *Client) RenderContextUsage() {

	bytes, err := json.Marshal(client.Session.Agent.ContextUsage)

	if err == nil {

		fmt.Fprintf(os.Stdout, "types.ContextUsage:%s\n", string(bytes))
		os.Stdout.Sync()

	}

}

func (client *Client) SetRole(role string) {

	if role == "user" || role == "assistant" {
		client.role = role
	}

}
