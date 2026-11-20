package main

import "fmt"
import "net/url"
import "os"
import "os/exec"

func CreateServer(listener_url *url.URL, model_name string, model_path string) (*Server, error) {

	_, err0 := exec.LookPath("llama-server")

	if err0 == nil {

		stat, err1 := os.Stat(model_path)

		if err1 == nil && stat.IsDir() == false {

			server := NewServer(listener_url, model_name, model_path)

			if server != nil {
				return server, nil
			} else {
				return nil, fmt.Errorf("test.CreateServer: Cannot create llama-server.")
			}

		} else {
			return nil, err1
		}

	} else {
		return nil, err0
	}

}
