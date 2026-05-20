package actions

import "exocomp/types"
import "exocomp/ui/web"
import "exocomp/ui/webview"
import "fmt"
import "os"
import "strings"

func Webview(agent *types.Agent, config *types.Config) {

	fmt.Fprintf(os.Stdout, "[config]:\n")
	fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", agent.Name, agent.Role, agent.Model, agent.Temperature)
	fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
	fmt.Fprintf(os.Stdout, "| Tools:   %s\n", strings.Join(agent.AllowedTools, ", "))
	fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
	fmt.Fprintf(os.Stdout, "\n")
	os.Stdout.Sync()

	server := web.NewServer(agent, config)
	client := webview.NewClient(server.URL)

	go client.Init()

	result := server.Init()

	if result == true {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

}
