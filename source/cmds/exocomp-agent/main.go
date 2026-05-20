package main

import "exocomp/actions"
import "exocomp/agents"
import "exocomp/types"
import "exocomp/utils/cli"
import "encoding/json"
import "fmt"
import "os"
import "strings"

func main() {

	var config *types.Config = nil
	var mode   string        = ""

	if len(os.Args) > 1 {

		tmp1 := strings.TrimSpace(os.Args[1])

		switch tmp1 {
		case "agent":
			mode = "agent"
		case "terminal":
			mode = "terminal"
		}

		config = cli.ParseConfig(os.Args[1:])

	}

	if mode != "" && config != nil {

		agent := agents.NewAgent(config)

		config.Update(
			agent.Name,
			agent.Role,
			agent.Model,
			config.Prompt,
			agent.Temperature,
		)

		err0 := os.MkdirAll(config.Sandbox, 0755)

		if err0 == nil {

			if config.Debug == true {
				tmp, _ := json.MarshalIndent(config, "", "\t")
				os.WriteFile(config.Sandbox + "/.exocomp-config.json", tmp, 0666)
			}

			switch mode {
			case "agent":

				actions.Agent(agent, config)

			case "terminal":

				actions.Terminal(agent, config, "user")

			default:

				actions.Usage([]string{"agent", "terminal"})

			}

		} else {

			fmt.Fprintf(os.Stderr, "Error: %s", err0.Error())
			os.Exit(1)

		}

	} else {

		actions.Usage([]string{"agent", "terminal"})

	}

}
