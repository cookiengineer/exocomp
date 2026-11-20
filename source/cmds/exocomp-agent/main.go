package main

import "exocomp/actions"
import "exocomp/agents"
import "exocomp/ui/jsonl"
import "exocomp/utils/cli"
import "fmt"
import "os"
import "strings"

func main() {

	arguments := os.Args[1:]
	show_usage := false

	if len(arguments) == 0 {
		show_usage = true
	}

	for _, argument := range arguments {

		tmp := strings.TrimSpace(argument)

		if strings.HasPrefix(tmp, "--") == false || tmp == "--help" {

			show_usage = true
			break

		}

	}

	if show_usage == false {

		config := cli.ParseConfig(arguments)

		if config != nil {

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

				jsonl.Run(agent, config)

			} else {

				fmt.Fprintf(os.Stderr, "Error: %s", err0.Error())
				os.Exit(1)

			}

		} else {

			actions.UsageAgent()

		}

	} else {

		actions.UsageAgent()

	}

}
