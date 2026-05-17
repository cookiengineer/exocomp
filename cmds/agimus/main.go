package main

import "exocomp/actions"
import "exocomp/agents"
import "exocomp/types"
import "exocomp/utils/cli"
import "fmt"
import "os"

func main() {

	var config *types.Config = nil

	if len(os.Args) > 1 {
		config = cli.ParseConfig(os.Args[1:])
	} else {
		config = cli.ParseConfig([]string{"planner"})
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

			actions.Terminal(agent, config, "assistant")

		} else {

			fmt.Fprintf(os.Stderr, "Error: %s", err0.Error())
			os.Exit(1)

		}

	} else {

		fmt.Fprintf(os.Stderr, "Error: %s", "Invalid \"agent\" flag")
		os.Exit(1)

	}

}
