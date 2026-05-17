//go:build !cgo

package webview

import "exocomp/utils/cli"
import "fmt"
import net_url "net/url"
import "os"
import "os/exec"
import "os/signal"
import "path/filepath"
import "syscall"

func getUserCacheDir() string {

	cache, err := os.UserCacheDir()

	if err == nil {
		return filepath.Join(cache, "exocomp")
	} else {
		return filepath.Join(os.TempDir(), "exocomp-cache")
	}

}

func getUserDataDir() string {

	config, err := os.UserConfigDir()

	if err == nil {
		return filepath.Join(config, "exocomp")
	} else {
		return filepath.Join(os.TempDir(), "exocomp")
	}

}

type Client struct {
	URL     *net_url.URL
	Process *exec.Cmd
}

func NewClient(url *net_url.URL) *Client {

	return &Client{
		URL: url,
	}

}

func (client *Client) Init() {

	browser, err0 := cli.FindBrowser()

	if err0 == nil {

		signals := make(chan os.Signal, 1)

		signal.Notify(
			signals,
			syscall.SIGINT,
			syscall.SIGTERM,
		)

		client.Process = exec.Command(
			browser,
			fmt.Sprintf("--app=%s", client.URL.String()),
			fmt.Sprintf("--user-data-dir=%s", getUserDataDir()),
			fmt.Sprintf("--disk-cache-dir=%s", getUserCacheDir()),
		)

		client.Process.Stdout = os.Stdout
		client.Process.Stderr = os.Stderr

		err1 := client.Process.Start()

		if err1 == nil {

			go func() {

				err := client.Process.Wait()

				if err != nil {
					fmt.Fprintf(os.Stderr, "Chromium exited with error: %s\n", err.Error())
				}

				os.Exit(0)

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

		} else {

			fmt.Fprintf(os.Stderr, "Error: %s\n", err1.Error())
			os.Exit(1)

		}

	} else {

		fmt.Fprintf(os.Stderr, "Error: %s\n", err0.Error())
		os.Exit(1)

	}

}

func (client *Client) Destroy() {

	if client.Process != nil && client.Process.Process != nil {

		err1 := client.Process.Process.Kill()

		if err1 != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err1.Error())
		}

		_, err2 := client.Process.Process.Wait()

		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err2.Error())
		}

	}

}
