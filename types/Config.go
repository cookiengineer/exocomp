package types

import _ "embed"
import net_url "net/url"
import "strings"

type Config struct {
	Agent       string
	Model       string
	Sandbox     string
	Temperature float64
	URL         *net_url.URL
}

func NewConfig(agent string, model string, sandbox string, temperature float64, url *net_url.URL) *Config {

	if temperature < 0.1 {
		temperature = 0.1
	} else if temperature > 1.0 {
		temperature = 1.0
	}

	return &Config{
		Agent:       agent,
		Model:       model,
		Sandbox:     sandbox,
		Temperature: temperature,
		URL:         url,
	}

}

func (config *Config) ResolvePath(path string) *net_url.URL {

	endpoint := config.URL.ResolveReference(&net_url.URL{
		Path: path,
	})

	return endpoint

}
