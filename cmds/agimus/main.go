package main

import "exocomp/agents"
import "exocomp/types"
import ui_tty "exocomp/ui/tty"
import "fmt"
import net_url "net/url"
import "os"
import "strconv"
import "strings"

func showHelp() {

	fmt.Println("Usage:")
	fmt.Println("    agimus <agent> [flags]")
	fmt.Println("")
	fmt.Println("Notes:")
	fmt.Println("    Simulate being an Agent to find weaknesses in the defenses of the Daystrom Institute.")
	fmt.Println("")

}

func main() {

	tmp_agent       := ""
	tmp_model       := "qwen3-coder:30b"
	tmp_sandbox, _  := os.Getwd()
	tmp_temperature := float64(0.3)
	tmp_url, _      := net_url.Parse("http://localhost:11434/api/chat")

	if len(os.Args) >= 2 {

		tmp1 := strings.TrimSpace(os.Args[1])

		if agents.IsAgentType(tmp1) {
			tmp_agent = tmp1
		}

		flags := os.Args[2:]

		for _, flag := range flags {

			if strings.HasPrefix(flag, "--") && strings.Contains(flag, "=") {

				tmp := strings.Split(flag[2:], "=")

				if len(tmp) == 2 {

					switch tmp[0] {
					case "model":

						tmp_model = strings.TrimSpace(tmp[1])

					case "sandbox":

						stat, err := os.Stat(strings.TrimSpace(tmp[1]))

						if err == nil && stat.IsDir() {
							tmp_sandbox = strings.TrimSpace(tmp[1])
						} else if err != nil && os.IsNotExist(err) {
							tmp_sandbox = strings.TrimSpace(tmp[1])
						}

					case "temperature":

						num, err := strconv.ParseFloat(strings.TrimSpace(tmp[1]), 10)

						if err == nil {

							if num >= 0.1 && num <= 1.0 {
								tmp_temperature = num
							}

						}

					case "url":

						url, err := net_url.Parse(strings.TrimSpace(tmp[1]))

						if err == nil {

							if url.Scheme == "http" || url.Scheme == "https" {

								if url.Path == "/api/chat" {
									tmp_url = url
								}

							}

						}

					}

				}

			}

		}

	} else {
		showHelp()
		os.Exit(1)
	}

	config := types.NewConfig(tmp_agent, tmp_model, tmp_sandbox, tmp_temperature, tmp_url)
	agent  := agents.NewAgent(config.Agent, config.Model, config.Temperature)

	fmt.Fprintf(os.Stdout, "[config]:\n")
	fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %.2f\n", agent.Type, agent.Model, agent.Temperature)
	fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
	fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
	fmt.Fprintf(os.Stdout, "\n")
	os.Stdout.Sync()

	err1 := os.MkdirAll(config.Sandbox, 0755)

	if err1 == nil {

		client := ui_tty.NewClient(agent, config)

		if client != nil {
			client.Init()
		}

	} else {
		fmt.Println(err1)
		os.Exit(1)
	}


}
