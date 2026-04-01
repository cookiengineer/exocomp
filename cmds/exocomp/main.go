package main

import "exocomp/agents"
import "exocomp/types"
import ui_tty "exocomp/ui/tty"
import ui_web "exocomp/ui/web"
import "fmt"
import net_url "net/url"
import "os"
import "strconv"
import "strings"

func showHelp() {

	fmt.Println("Usage:")
	fmt.Println("    exocomp <ui> <agent> [flags]")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("    <ui> string            UI type")
	fmt.Println("                           Either of: tty, web")
	fmt.Println("")
	fmt.Println("    <agent> string         LLM agent type")
	fmt.Println("                           Either of: architect, coder, manager, tester")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("    --model string         LLM model identifier (ollama format)")
	fmt.Println("                           Run \"ollama list\" to see available models")
	fmt.Println("                           Examples: qwen3-coder:30b, codestral:22b")
	fmt.Println("                           (default: \"qwen3-coder:30b\")")
	fmt.Println("")
	fmt.Println("    --temperature float    LLM sampling temperature (0.1-1.0)")
	fmt.Println("                           Lower = more deterministic, fewer hallucinations")
	fmt.Println("                           Higher = more creative, more hallucinations")
	fmt.Println("                           (default: 0.3)")
	fmt.Println("")
	fmt.Println("    --sandbox string       Path to sandbox directory")
	fmt.Println("                           (default: current working directory)")
	fmt.Println("")
	fmt.Println("    --url string           API endpoint for LLM backend")
	fmt.Println("                           (default: \"http://localhost:11434/api/chat\")")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("    exocomp tty coder --model=\"codestral:22b\" --temperature=0.1;")
	fmt.Println("    exocomp tty architect --model=\"qwen3.5:35b\" --temperature=0.7;")
	fmt.Println("")
	fmt.Println("    exocomp web architect --model=\"qwen3.5:35b\" --temperature=0.7;")
	fmt.Println("")

}

func main() {

	tmp_ui          := ""
	tmp_agent       := ""
	tmp_model       := "qwen3-coder:30b"
	tmp_sandbox, _  := os.Getwd()
	tmp_temperature := float64(0.3)
	tmp_url, _      := net_url.Parse("http://localhost:11434/api/chat")

	if len(os.Args) >= 2 {

		tmp1 := strings.TrimSpace(os.Args[1])
		tmp2 := strings.TrimSpace(os.Args[2])

		if (tmp1 == "tty" || tmp1 == "web") && agents.IsAgentType(tmp2) {

			tmp_ui    = tmp1
			tmp_agent = tmp2

		} else {

			showHelp()
			os.Exit(1)

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

		if tmp_ui == "tty" {

			client := ui_tty.NewClient(agent, config)
			client.Init()

		} else if tmp_ui == "web" {

			server := ui_web.NewServer(agent, config)
			client := ui_web.NewClient(server.URL)

			go client.Init()
			server.Init()

			// TODO: client.Destroy()

		} else {

			showHelp()
			os.Exit(1)

		}

	} else {
		fmt.Println(err1)
		os.Exit(1)
	}

}
