package main

import "exocomp/agents"
import "exocomp/types"
import utils_fmt "exocomp/utils/fmt"
import ui_jsonl "exocomp/ui/jsonl"
import ui_tty "exocomp/ui/tty"
import ui_web "exocomp/ui/web"
import ui_webview "exocomp/ui/webview"
import "fmt"
import net_url "net/url"
import "os"
import "strconv"
import "strings"

func showHelp() {

	fmt.Println("Usage:")
	fmt.Println("    exocomp <ui> [flags]")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("    <ui> string            UI type")
	fmt.Println("                           Either of: jsonl, tty, web, webview")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("")
	fmt.Println("    --name string          LLM agent name")
	fmt.Println("                           (default: \"Peanut Hamper\")")
	fmt.Println("")
	fmt.Println("    --agent string         LLM agent type")
	fmt.Println("                           Either of: architect, coder, manager, summarizer, tester")
	fmt.Println("")
	fmt.Println("    --model string         LLM agent model (ollama format)")
	fmt.Println("                           Run \"ollama list\" to see available models")
	fmt.Println("                           Examples: qwen3-coder:30b, codestral:22b")
	fmt.Println("                           (default: \"qwen3-coder:30b\")")
	fmt.Println("")
	fmt.Println("    --temperature float    LLM agent sampling temperature (0.1-1.0)")
	fmt.Println("                           Lower = more deterministic, fewer hallucinations")
	fmt.Println("                           Higher = more creative, more hallucinations")
	fmt.Println("                           (default: 0.3)")
	fmt.Println("")
	fmt.Println("    --prompt string        Initial LLM instructions prompt")
	fmt.Println("                           (default: \"\")")
	fmt.Println("")
	fmt.Println("    --sandbox string       Path to sandbox directory")
	fmt.Println("                           (default: current working directory)")
	fmt.Println("")
	fmt.Println("    --url string           API endpoint for LLM backend")
	fmt.Println("                           (default: \"http://localhost:11434/api/chat\")")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("")
	fmt.Println("    # single-agent mode")
	fmt.Println("    exocomp tty --agent=architect")
	fmt.Println("    exocomp web --agent=architect --model=\"qwen3.5:35b\" --temperature=\"0.7\";")
	fmt.Println("")
	fmt.Println("    # multi-agent mode")
	fmt.Println("    exocomp tty --agent=manager")
	fmt.Println("    exocomp web --agent=manager --model=\"codestral:22b\" --temperature=\"0.2\";")
	fmt.Println("")

}

func main() {

	tmp_ui            := ""
	tmp_name          := ""
	tmp_agent         := ""
	tmp_debug         := false
	tmp_model         := "qwen3-coder:30b"
	tmp_playground, _ := os.Getwd()
	tmp_prompt        := ""
	tmp_sandbox, _    := os.Getwd()
	tmp_temperature   := float64(0.3)
	tmp_url, _        := net_url.Parse("http://localhost:11434/api/chat")

	if len(os.Args) >= 2 {

		tmp1 := strings.TrimSpace(os.Args[1])

		if (tmp1 == "jsonl" || tmp1 == "tty" || tmp1 == "web" || tmp1 == "webview") {

			tmp_ui = tmp1

		} else {

			showHelp()
			os.Exit(1)

		}

		flags := os.Args[1:]

		for _, flag := range flags {

			if strings.HasPrefix(flag, "--") && strings.Contains(flag, "=") {

				tmp := strings.Split(flag[2:], "=")

				if len(tmp) == 2 {

					switch tmp[0] {
					case "name":

						if agents.IsAgentName(tmp[1]) {
							tmp_name = utils_fmt.FormatAgentName(tmp[1])
						}

					case "agent":

						if agents.IsAgentType(tmp[1]) {
							tmp_agent = strings.TrimSpace(tmp[1])
						}

					case "model":

						tmp_model = strings.TrimSpace(tmp[1])

					case "playground":

						stat, err := os.Stat(strings.TrimSpace(tmp[1]))

						if err == nil && stat.IsDir() {
							tmp_playground = strings.TrimSpace(tmp[1])
						} else if err != nil && os.IsNotExist(err) {
							tmp_playground = strings.TrimSpace(tmp[1])
						}

					case "prompt":

						tmp_prompt = utils_fmt.FormatSingleLine(tmp[1])

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

			} else if strings.HasPrefix(flag, "--") {

				tmp := strings.TrimSpace(flag[2:])

				switch tmp {
				case "debug":
					tmp_debug = true
				}

			}

		}

	} else {
		showHelp()
		os.Exit(1)
	}

	if tmp_playground == tmp_sandbox || strings.HasPrefix(tmp_sandbox, tmp_playground + string(os.PathSeparator)) {

		config := types.NewConfig(tmp_name, tmp_agent, tmp_debug, tmp_model, tmp_playground, tmp_prompt, tmp_sandbox, tmp_temperature, tmp_url)
		agent  := agents.NewAgent(config.Name, config.Agent, config.Model, config.Temperature)

		err1 := os.MkdirAll(config.Sandbox, 0755)

		if err1 == nil {

			if tmp_ui == "jsonl" {

				os.Stdout.Sync()

				client := ui_jsonl.NewClient(agent, config)
				client.Init()

			} else if tmp_ui == "tty" {

				fmt.Fprintf(os.Stdout, "[config]:\n")
				fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", agent.Name, agent.Type, agent.Model, agent.Temperature)
				fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
				fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
				fmt.Fprintf(os.Stdout, "\n")
				os.Stdout.Sync()

				client := ui_tty.NewClient(agent, config)
				client.Init()

			} else if tmp_ui == "web" {

				server := ui_web.NewServer(agent, config)

				fmt.Fprintf(os.Stdout, "[config]:\n")
				fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", agent.Name, agent.Type, agent.Model, agent.Temperature)
				fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
				fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
				fmt.Fprintf(os.Stdout, "| Web:     %s\n", server.URL.String())
				fmt.Fprintf(os.Stdout, "\n")
				os.Stdout.Sync()

				server.Init()

			} else if tmp_ui == "webview" {

				fmt.Fprintf(os.Stdout, "[config]:\n")
				fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", agent.Name, agent.Type, agent.Model, agent.Temperature)
				fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", config.Sandbox)
				fmt.Fprintf(os.Stdout, "| URL:     %s\n", config.URL.String())
				fmt.Fprintf(os.Stdout, "\n")
				os.Stdout.Sync()

				server := ui_web.NewServer(agent, config)
				client := ui_webview.NewClient(server.URL)

				go client.Init()
				server.Init()

			} else {

				showHelp()
				os.Exit(1)

			}

		} else {
			fmt.Println(err1)
			os.Exit(1)
		}

	} else {

		fmt.Println("Invalid sandbox path, must be inside playground.")
		os.Exit(1)

	}

}
