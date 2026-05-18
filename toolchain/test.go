package main

import "exocomp-toolchain/utils"
import "fmt"
import "net/url"
import "os"
import "strings"
import "time"

func main() {

	base_dir, err0 := utils.GetRoot()
	base_url, err1 := url.Parse("http://localhost:11434/v1/models")

	if err0 == nil && err1 == nil {

		err2 := utils.CheckModel(base_url, "qwen3-coder:30b")

		fmt.Println(err2)

		if err2 == nil {

			// TODO: server is managed externally

		} else if strings.Contains(err2.Error(), "connection refused") {

			// TODO: kill process that runs on port 11434

			fmt.Fprint(os.Stdout, "== Start llama-server ==\n")

			server, err := utils.CreateServer(
				base_url,
				"qwen3-coder:30b",
				"/opt/llama/models/qwen3-coder-30b-a3b-instruct-q8_0.gguf",
			)

			if err == nil && server != nil {

				server.Start()

				// TODO: Wait 30 seconds?

				time.Sleep(30 * time.Second)

				fmt.Fprint(os.Stdout, "\n")
				fmt.Fprint(os.Stdout, "== Test tools ==\n")

				fmt.Fprintf(os.Stdout, "base_dir: %s\n", base_dir)

				// TODO

				fmt.Fprint(os.Stdout, "\n")
				fmt.Fprint(os.Stdout, "== Test agents ==\n")

				server.Stop()

			} else {
				fmt.Fprint(os.Stderr, "Couldn't create llama.cpp server: %s", err.Error())
			}

		}

	} else {

		if err0 != nil {
			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err0.Error())
		}

		if err1 != nil {
			fmt.Fprintf(os.Stderr, "!! Error: %s\n", err1.Error())
		}

	}

}
