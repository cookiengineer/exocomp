package tty

import "exocomp/agents"
import "exocomp/ollama"
import "exocomp/types"
import "fmt"
import "os"
import "os/signal"
import "syscall"

type Client struct {
	Session  *ollama.Session
	Renderer *Renderer
}

func NewClient(agent *agents.Agent, config *types.Config) *Client {

	session, err0 := ollama.NewSession(agent, config)

	if err0 == nil {

		renderer := NewRenderer(session)

		return &Client{
			Session:  session,
			Renderer: renderer,
		}

	} else {
		return nil
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
		client.Renderer.InputLoop()
		signals<-syscall.SIGINT
	}()

	go client.Renderer.RenderLoop()

	select {
	case sig := <-signals:

		switch sig {
		case syscall.SIGINT:

			client.Renderer.Destroy()
			fmt.Fprintf(os.Stdout, "Received signal: %s\n", "SIGINT")
			os.Exit(0)

		case syscall.SIGTERM:

			client.Renderer.Destroy()
			fmt.Fprintf(os.Stdout, "Received signal: %s\n", "SIGTERM")
			os.Exit(0)

		default:

			client.Renderer.Destroy()
			fmt.Printf("Received signal: %s\n", sig.String())
			os.Exit(0)

		}

	}

}
