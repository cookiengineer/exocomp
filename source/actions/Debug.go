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

	if types.GlobalConfig != nil {

		fmt.Fprintf(os.Stderr, "[global config]:\n")
		fmt.Fprintf(os.Stderr, "| Name:        %s\n", types.GlobalConfig.Name)
		fmt.Fprintf(os.Stderr, "| Role:        %s\n", types.GlobalConfig.Role)
		fmt.Fprintf(os.Stderr, "| Model:       %s\n", types.GlobalConfig.Model)
		fmt.Fprintf(os.Stderr, "| Prompt:      %d bytes\n", len(types.GlobalConfig.Prompt))
		fmt.Fprintf(os.Stderr, "| Temperature: %.2f\n", types.GlobalConfig.Temperature)
		fmt.Fprintf(os.Stderr, "| Sandbox:     %s\n", types.GlobalConfig.Sandbox)

		if types.GlobalConfig.URL != nil {
			fmt.Fprintf(os.Stderr, "| URL:         %s\n", types.GlobalConfig.URL.String())
		} else {
			fmt.Fprintf(os.Stderr, "| URL:         %s\n", "")
		}

		fmt.Fprintf(os.Stderr, "| Debug:       %t\n", types.GlobalConfig.Debug)
		fmt.Fprintf(os.Stderr, "| Providers:\n")

		for model, provider := range types.GlobalConfig.Providers {

			if provider.URL != nil {

				if provider.Alias != "" {
					fmt.Fprintf(os.Stderr, "|> \"%s\" via \"%s\" as \"%s\"\n", model, provider.URL.String(), provider.Alias)
				} else {
					fmt.Fprintf(os.Stderr, "|> \"%s\" via \"%s\"\n", model, provider.URL.String())
				}

			}

		}

		fmt.Fprintf(os.Stderr, "\n")

	}

	fmt.Fprintf(os.Stderr, "[config]:\n")
	fmt.Fprintf(os.Stderr, "| Name:        %s\n", config.Name)
	fmt.Fprintf(os.Stderr, "| Role:        %s\n", config.Role)
	fmt.Fprintf(os.Stderr, "| Model:       %s\n", config.Model)
	fmt.Fprintf(os.Stderr, "| Prompt:      %d bytes\n", len(config.Prompt))
	fmt.Fprintf(os.Stderr, "| Temperature: %.2f\n", config.Temperature)
	fmt.Fprintf(os.Stderr, "| Sandbox:     %s\n", config.Sandbox)
	fmt.Fprintf(os.Stderr, "| URL:         %s\n", config.URL.String())
	fmt.Fprintf(os.Stderr, "| Debug:       %t\n", config.Debug)
	fmt.Fprintf(os.Stderr, "| Providers:\n")

	for model, provider := range config.Providers {

		if provider.Alias != "" {
			fmt.Fprintf(os.Stderr, "|> \"%s\" as \"%s\": \"%s\", \"%s\"\n", model, provider.Alias, provider.URL.String(), provider.Token)
		} else {
			fmt.Fprintf(os.Stderr, "|> \"%s\": \"%s\", \"%s\"\n", model, provider.URL.String(), provider.Token)
		}

	}

	fmt.Fprintf(os.Stderr, "\n")

	os.Stdout.Sync()
	os.Stderr.Sync()

	fmt.Fprintf(os.Stderr, "\n")

}
