package actions

import "exocomp/types"
import "exocomp/ui/jsonl"
import "fmt"
import "os"
import "strings"

func Agent(agent *types.Agent, config *types.Config) {

	fmt.Fprintf(os.Stderr, "[config]:\n")
	fmt.Fprintf(os.Stderr, "| Agent:    %s | %s | %s | %.2f\n", agent.Name, agent.Type, agent.Model, agent.Temperature)
	fmt.Fprintf(os.Stderr, "| Sandbox:  %s\n", config.Sandbox)
	fmt.Fprintf(os.Stderr, "| Tools:    %s\n", strings.Join(agent.AllowedTools, ", "))
	fmt.Fprintf(os.Stderr, "| URL:      %s\n", config.URL.String())
	fmt.Fprintf(os.Stderr, "\n")

	os.Stdout.Sync()
	os.Stderr.Sync()

	client := jsonl.NewClient(agent, config)
	client.SetRole("user")

	result := client.Init()

	if result == true {
		os.Exit(0)
	} else {
		os.Exit(1)
	}

}
