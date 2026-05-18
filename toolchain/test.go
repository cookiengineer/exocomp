package main

import "exocomp-toolchain/utils"

func main() {

	base_dir, err0 := utils.GetRoot()
	base_url, err1 := url.Parse("http://localhost:11434/v1/models")

	if err0 == nil && err1 == nil {

		err2 := utils.CheckServer(base_url, "qwen3-coder:30b")

		if err2 == nil {

			// TODO: server is managed externally

		} else {

			// TODO: kill process that runs on port 11434

			fmt.Fprint(os.Stdout, "== Start llama-server ==\n")

			err := utils.StartServer(base_url, "qwen3-coder:30b")

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Test tools ==\n")

			// TODO

			fmt.Fprint(os.Stdout, "\n")
			fmt.Fprint(os.Stdout, "== Test agents ==\n")

			// TODO

			// TODO: Start llama.cpp

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
