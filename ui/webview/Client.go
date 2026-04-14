package webview

import "fmt"
import net_url "net/url"
import "os"
import "os/signal"
import "syscall"
import "time"

type Client struct {
	URL     *net_url.URL
	Webview WebView
}

func NewClient(url *net_url.URL) *Client {

	return &Client{
		URL: url,
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

		time.Sleep(500 * time.Millisecond)

		client.Webview = New(true)

		client.Webview.SetTitle("Exocomp")
		client.Webview.SetSize(1024, 768, HintNone)
		client.Webview.Navigate(client.URL.String())
		client.Webview.Run()

	}()

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

	if client.Webview != nil {
		client.Webview.Destroy()
	}

}
