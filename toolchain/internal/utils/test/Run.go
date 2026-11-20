package test

import "fmt"
import "net/url"
import "os"
import "strings"
import "time"

func Run(base_dir string) error {

	base_url, err0 := url.Parse("http://localhost:11434/v1/models")

	if err0 == nil {

		err1 := CheckModel(base_url, "huihui_ai/Qwen3.8-abliterated:27b")

		if err1 == nil {
			return nil
		} else if strings.Contains(err1.Error(), "connection refused") == false {
			return err1
		} else {

			fmt.Fprint(os.Stdout, "== Start llama-server ==\n")

			server, err2 := CreateServer(
				base_url,
				"huihui_ai/Qwen3.8-abliterated:27b",
				"/opt/llama/models/Huihui-Qwen3.8-27B-abliterated.Q4_K_M.gguf",
			)

			if err2 == nil {

				err3 := server.Start()

				if err3 == nil {

					time.Sleep(30 * time.Second)

					fmt.Fprint(os.Stdout, "\n")
					fmt.Fprint(os.Stdout, "== Test tools ==\n")
					fmt.Fprintf(os.Stdout, "base_dir: %s\n", base_dir)
					fmt.Fprint(os.Stdout, "\n")
					fmt.Fprint(os.Stdout, "== Test agents ==\n")

					return server.Stop()

				} else {
					return err3
				}

			} else {
				return err2
			}

		}

	} else {
		return err0
	}

}
