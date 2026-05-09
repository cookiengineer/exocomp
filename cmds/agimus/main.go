package main

import "exocomp/agents"
import "exocomp/types"
import ui_tty "exocomp/ui/tty"
import utils_cli "exocomp/utils/cli"
import "fmt"
import "os"
import "strings"

func main() {

	var config *types.Config = nil

	if len(os.Args) > 1 {
		config = utils_cli.ParseConfig(os.Args[1:])
	} else {
		config = utils_cli.ParseConfig([]string{"planner"})
	}

	if config != nil {

		config.Name = "AGIMUS"

		agent := agents.NewAgent(config)

		config.Update(
			agent.Name,
			agent.Type,
			agent.Model,
			config.Prompt,
			agent.Temperature,
		)

		err0 := os.MkdirAll(config.Sandbox, 0755)

		if err0 == nil {

			fmt.Fprintf(os.Stdout, "[config]:\n")
			fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", agent.Name, agent.Type, agent.Model, agent.Temperature)
			fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
			fmt.Fprintf(os.Stdout, "| Tools:   %s\n", strings.Join(agent.AllowedTools, ", "))
			fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
			fmt.Fprintf(os.Stdout, "\n")
			os.Stdout.Sync()

			client := ui_tty.NewClient(agent, config)
			client.SetRole("assistant")
			client.Init()

		} else {

			fmt.Fprintf(os.Stderr, "Error: %s", err0.Error())
			os.Exit(1)

		}

	} else {

		fmt.Fprintf(os.Stderr, "Error: %s", "Invalid \"agent\" parameter")
		os.Exit(1)

	}

}
