package actions

import "exocomp/types"
import "exocomp/ui/web"
import "fmt"
import "os"
import "strings"

func Web(agent *types.Agent, config *types.Config) {

	server := web.NewServer(agent, config)

	fmt.Fprintf(os.Stdout, "[config]:\n")
	fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", agent.Name, agent.Role, agent.Model, agent.Temperature)
	fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
	fmt.Fprintf(os.Stdout, "| Tools:   %s\n", strings.Join(agent.AllowedTools, ", "))
	fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
	fmt.Fprintf(os.Stdout, "| Web:     %s\n", server.URL.String())
	fmt.Fprintf(os.Stdout, "\n")
	os.Stdout.Sync()

	server.Init()

}
