package main

import "exocomp-toolchain/utils"
import "fmt"
import "net/url"
import "os"
import "strings"

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

			// err3 := utils.StartServer(base_url, "qwen3-coder:30b")

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Test tools ==\n")

			// TODO

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Test agents ==\n")

			// TODO

			// TODO: Start llama.cpp

			if 1 == 2 {
				fmt.Println(base_dir)
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
