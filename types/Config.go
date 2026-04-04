package types

import _ "embed"
import net_url "net/url"
import "os"

type Config struct {
	Name        string
	Agent       string
	Model       string
	Playground  string
	Sandbox     string
	Temperature float64
	URL         *net_url.URL
}

func NewConfig(name string, agent string, model string, sandbox string, temperature float64, url *net_url.URL) *Config {

	if temperature < 0.1 {
		temperature = 0.1
	} else if temperature > 1.0 {
		temperature = 1.0
	}

	base := os.TempDir()
	playground, err := os.MkdirTemp(base, "exocomp-playground-*")

	if err == nil {

		return &Config{
			Name:        name,
			Agent:       agent,
			Model:       model,
			Playground:  playground,
			Sandbox:     sandbox,
			Temperature: temperature,
			URL:         url,
		}

	} else {
		panic(err)
	}

}

func (config *Config) ResolvePath(path string) *net_url.URL {

	endpoint := config.URL.ResolveReference(&net_url.URL{
		Path: path,
	})

	return endpoint

}
