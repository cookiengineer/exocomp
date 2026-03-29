package main

import "exocomp/agents"
import "exocomp/ollama"
import "exocomp/types"
import ui_tty "exocomp/ui/tty"
import "fmt"
import net_url "net/url"
import "os"
import "os/signal"
import "strconv"
import "strings"
import "syscall"

func showHelp() {

	fmt.Println("Usage:")
	fmt.Println("    exocomp <agent> [flags]")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("    <agent> string         Type of agent")
	fmt.Println("                           Either of: manager, coder, tester")
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
	fmt.Println("    exocomp coder --model=\"codestral:22b\" --temperature=0.1;")
	fmt.Println("    exocomp manager --model=\"qwen3.5:35b\" --temperature=0.7;")
	fmt.Println("")

}

func main() {

	tmp_agent       := ""
	tmp_model       := "qwen3-coder:30b"
	tmp_sandbox, _  := os.Getwd()
	tmp_temperature := float64(0.3)
	tmp_url, _      := net_url.Parse("http://localhost:11434/api/chat")

	if len(os.Args) >= 2 {

		tmp := strings.TrimSpace(os.Args[1])

		if agents.IsAgentType(tmp) == true {

			tmp_agent = tmp

		} else {

			showHelp()
			os.Exit(1)

		}

		if len(os.Args) >= 3 {

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

		session, err2 := ollama.NewSession(agent, config)

		if err2 == nil {

			renderer := ui_tty.NewRenderer(session)
			signals  := make(chan os.Signal, 1)

			signal.Notify(
				signals,
				syscall.SIGINT,
				syscall.SIGTERM,
			)

			go func() {
				renderer.InputLoop()
				signals<-syscall.SIGINT
			}()

			go renderer.RenderLoop()

			select {
			case sig := <-signals:

				switch sig {
				case syscall.SIGINT:

					renderer.Destroy()
					fmt.Println("Received signal: SIGINT")
					os.Exit(0)

				case syscall.SIGTERM:

					renderer.Destroy()
					fmt.Println("Received signal: SIGTERM")
					os.Exit(0)

				default:

					renderer.Destroy()
					fmt.Printf("Received signal: %s\n", sig.String())
					os.Exit(0)

				}

			}

		} else {
			fmt.Println(err2)
			os.Exit(1)
		}

	} else {
		fmt.Println(err1)
		os.Exit(1)
	}

}
