package tty

import "exocomp/agents"
import "exocomp/schemas"
import "exocomp/types"
import "bufio"
import "encoding/json"
import "fmt"
import "os"
import "os/signal"
import "strings"
import "syscall"

type Client struct {
	Session  *types.Session
	Renderer *Renderer
}

func NewClient(agent *agents.Agent, config *types.Config) *Client {

	session  := types.NewSession(agent, config)
	renderer := NewRenderer(session)

	return &Client{
		Session:  session,
		Renderer: renderer,
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

		prompt := strings.TrimSpace(scanner.Text())

		if prompt != "" && client.Session != nil {

			go func() {

				err := client.Session.SendChatRequest(schemas.Message{
					Role:    "user",
					Content: prompt,
				})

				if err != nil {

					if client.Session.Config.Debug == true {

						bytes1, _ := json.MarshalIndent(client.Session.Config, "", "\t")
						bytes2, _ := json.MarshalIndent(schemas.ChatRequest{
							Model:       client.Session.Agent.Model,
							Temperature: client.Session.Agent.Temperature,
							Messages:    client.Session.Messages,
							Stream:      false,
							Tools:       client.Session.Tools,
							ToolChoice:  "auto",
						}, "", "\t")

						os.WriteFile(client.Session.Config.Sandbox + "/.exocomp-debug-config.json", bytes1, 0666)
						os.WriteFile(client.Session.Config.Sandbox + "/.exocomp-debug-chatrequest.json", bytes2, 0666)

						fmt.Fprintf(os.Stderr, "\nFatal Error: %s\n", err.Error())
						os.Exit(1)

					}

				}

			}()

		}

	}

}

func (client *Client) Destroy() {

	if client.Renderer != nil {
		client.Renderer.Destroy()
	}

}
