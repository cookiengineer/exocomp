package main

import "bufio"
import "fmt"
import "os"
import "exocomp/config"
import "exocomp/ollama"

func main() {

	config, err0 := config.ParseConfig()

	if err0 == nil {

		if config.Verbose == true {
			fmt.Println("Model:   %s", config.Model)
			fmt.Println("URL:     %s", config.URL.String())
			fmt.Println("Sandbox: %s", config.Sandbox)
			fmt.Println("Verbose: %b", config.Verbose)
		}

		err1 := os.MkdirAll(config.Sandbox, 0755)

		if err1 == nil {

			session, err2 := ollama.NewSession(config)

			if err2 == nil {

				scanner := bufio.NewScanner(os.Stdin)

				for {

					fmt.Printf("%s > ", config.Model)

					if scanner.Scan() == true {

						prompt := scanner.Text()

						if prompt == "exit" || prompt == "quit" {

							break

						} else {

							response, err3 := session.Query(prompt)

							if err3 == nil {
								fmt.Printf("%s\n", response)
							}

						}

					} else {
						break
					}

				}

				os.Exit(0)

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
