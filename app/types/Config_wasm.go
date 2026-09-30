//go:build wasm

package types

import "github.com/cookiengineer/gooey/components/app"
import "encoding/json"

// BootstrapConfig mirrors the old frontend `types/Config.mjs` BootstrapConfig
// helper. It blocks until the request finished; callers must run it inside a
// goroutine.
func BootstrapConfig(client *app.Client, agent string) (*Config, error) {

	path := "/api/session/config"

	if agent != "" {
		path = "/api/session/config/" + agent
	}

	response, err0 := client.Read(path)

	if err0 != nil {
		return nil, err0
	}

	config := NewConfig()
	err1 := json.Unmarshal(response.Body, config)

	if err1 != nil {
		return nil, err1
	}

	return config, nil

}
