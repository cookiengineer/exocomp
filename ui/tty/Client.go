package tty

import "exocomp/agents"
import "exocomp/types"
import "fmt"
import "os"
import "os/signal"
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
		client.Renderer.InputLoop()
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

func (client *Client) Destroy() {

	if client.Renderer != nil {
		client.Renderer.Destroy()
	}

}
