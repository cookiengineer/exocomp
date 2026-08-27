package actions

import "exocomp/types"
import "fmt"
import "os"
import "strings"

func Debug(agent *types.Agent, config *types.Config) {

	fmt.Fprintf(os.Stderr, "[agent]:\n")
	fmt.Fprintf(os.Stderr, "| Agent: %s | %s | %s | %.2f\n", agent.Name, agent.Role, agent.Model, agent.Temperature)
	fmt.Fprintf(os.Stderr, "| Tools: %s\n", strings.Join(agent.AllowedTools, ", "))
	fmt.Fprintf(os.Stderr, "\n")

	fmt.Fprintf(os.Stderr, "[config]:\n")
	fmt.Fprintf(os.Stderr, "| Name:        %s\n", config.Name)
	fmt.Fprintf(os.Stderr, "| Role:        %s\n", config.Role)
	fmt.Fprintf(os.Stderr, "| Model:       %s\n", config.Model)
	fmt.Fprintf(os.Stderr, "| Prompt:      %d bytes\n", len(config.Prompt))
	fmt.Fprintf(os.Stderr, "| Temperature: %.2f\n", config.Temperature)
	fmt.Fprintf(os.Stderr, "| Sandbox:     %s\n", config.Sandbox)
	fmt.Fprintf(os.Stderr, "| URL:         %s\n", config.URL.String())
	fmt.Fprintf(os.Stderr, "| Debug:       %v\n", config.Debug)
	fmt.Fprintf(os.Stderr, "| Tokens:\n")

	for model, token := range config.Tokens {
		fmt.Fprintf(os.Stderr, "|> \"%s\": \"%s\"\n", model, token)
	}

	fmt.Fprintf(os.Stderr, "\n")

	os.Stdout.Sync()
	os.Stderr.Sync()

	fmt.Fprintf(os.Stderr, "\n")

}
