package main

import "exocomp/agents"
import "exocomp/types"
import ui_jsonl "exocomp/ui/jsonl"
import ui_tty "exocomp/ui/tty"
import ui_web "exocomp/ui/web"
import ui_webview "exocomp/ui/webview"
import utils_cli "exocomp/utils/cli"
import "encoding/json"
import "fmt"
import "os"
import "strings"

func main() {

	var config *types.Config = nil
	var mode   string        = ""

	if len(os.Args) > 1 {

		tmp1 := strings.TrimSpace(os.Args[1])

		if (tmp1 == "jsonl" || tmp1 == "tty" || tmp1 == "web" || tmp1 == "webview") {
			mode = tmp1
		}

		config = utils_cli.ParseConfig(os.Args[1:])

	}

	if mode != "" && config != nil {

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

			if config.Debug == true {
				tmp, _ := json.MarshalIndent(config, "", "\t")
				os.WriteFile(config.Sandbox + "/.exocomp-config.json", tmp, 0666)
			}

			switch mode {
			case "jsonl":

				fmt.Fprintf(os.Stderr, "[config]:\n")
				fmt.Fprintf(os.Stderr, "| Agent:    %s | %s | %s | %.2f\n", agent.Name, agent.Type, agent.Model, agent.Temperature)
				fmt.Fprintf(os.Stderr, "| Sandbox:  %s\n", config.Sandbox)
				fmt.Fprintf(os.Stderr, "| Tools:    %s\n", strings.Join(agent.AllowedTools, ", "))
				fmt.Fprintf(os.Stderr, "| URL:      %s\n", config.URL.String())
				fmt.Fprintf(os.Stderr, "\n")

				os.Stdout.Sync()
				os.Stderr.Sync()

				client := ui_jsonl.NewClient(agent, config)
				client.SetRole("user")
				client.Init()

			case "tty":

				fmt.Fprintf(os.Stdout, "[config]:\n")
				fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", agent.Name, agent.Type, agent.Model, agent.Temperature)
				fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
				fmt.Fprintf(os.Stdout, "| Tools:   %s\n", strings.Join(agent.AllowedTools, ", "))
				fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
				fmt.Fprintf(os.Stdout, "\n")
				os.Stdout.Sync()

				client := ui_tty.NewClient(agent, config)
				client.SetRole("user")
				client.Init()

			case "web":

				server := ui_web.NewServer(agent, config)

				fmt.Fprintf(os.Stdout, "[config]:\n")
				fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", agent.Name, agent.Type, agent.Model, agent.Temperature)
				fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
				fmt.Fprintf(os.Stdout, "| Tools:   %s\n", strings.Join(agent.AllowedTools, ", "))
				fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
				fmt.Fprintf(os.Stdout, "| Web:     %s\n", server.URL.String())
				fmt.Fprintf(os.Stdout, "\n")
				os.Stdout.Sync()

				server.Init()

			case "webview":

				fmt.Fprintf(os.Stdout, "[config]:\n")
				fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", agent.Name, agent.Type, agent.Model, agent.Temperature)
				fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
				fmt.Fprintf(os.Stdout, "| Tools:   %s\n", strings.Join(agent.AllowedTools, ", "))
				fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
				fmt.Fprintf(os.Stdout, "\n")
				os.Stdout.Sync()

				server := ui_web.NewServer(agent, config)
				client := ui_webview.NewClient(server.URL)

				go client.Init()
				server.Init()

			default:

				utils_cli.PrintUsage()
				os.Exit(1)

			}

		} else {

			fmt.Fprintf(os.Stderr, "Error: %s", err0.Error())
			os.Exit(1)

		}

	} else {

		utils_cli.PrintUsage()
		os.Exit(1)

	}

}
