package actions

import "exocomp/types"
import "exocomp/ui/tty"
import "fmt"
import "os"
import "strings"

func Terminal(agent *types.Agent, config *types.Config, role string) {

	fmt.Fprintf(os.Stdout, "[config]:\n")
	fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", agent.Name, agent.Type, agent.Model, agent.Temperature)
	fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
	fmt.Fprintf(os.Stdout, "| Tools:   %s\n", strings.Join(agent.AllowedTools, ", "))
	fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
	fmt.Fprintf(os.Stdout, "\n")
	os.Stdout.Sync()

	client := tty.NewClient(agent, config)

	if role != "" {
		client.SetRole(role)
	}

	result := client.Init()

	if result == true {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

}
