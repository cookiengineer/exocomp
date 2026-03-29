package web

import webview "github.com/webview/webview_go"
import net_url "net/url"
import "time"

type Client struct {
	URL     *net_url.URL
	Webview webview.WebView
}

func NewClient(url *net_url.URL) *Client {

	return &Client{
		URL: url,
	}

}

func (client *Client) Init() {

	time.Sleep(500 * time.Millisecond)

	client.Webview = webview.New(true)

	client.Webview.SetTitle("Exocomp")
	client.Webview.SetSize(1024, 768, webview.HintNone)
	client.Webview.Navigate(client.URL.String())
	client.Webview.Run()

}

func (client *Client) Destroy() {

	if client.Webview != nil {
		client.Webview.Destroy()
	}

}
