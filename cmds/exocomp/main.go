package main

import "exocomp/agents"
import "exocomp/config"
import "exocomp/ollama"
import "fmt"
import "os"
import "os/signal"
import "syscall"

func main() {

	config, err0 := config.ParseConfig()

	if err0 == nil {

		agent := agents.NewAgent(config.Agent)

		if config.Verbose == true {
			fmt.Println("Model:   %s", config.Model)
			fmt.Println("URL:     %s", config.URL.String())
			fmt.Println("Sandbox: %s", config.Sandbox)
			fmt.Println("Verbose: %b", config.Verbose)
		}

		err1 := os.MkdirAll(config.Sandbox, 0755)

		if err1 == nil {

			session, err2 := ollama.NewSession(agent, config)

			if err2 == nil {

				renderer := ollama.NewRenderer(session)
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

	} else {
		fmt.Println(err0)
		os.Exit(1)
	}

}
