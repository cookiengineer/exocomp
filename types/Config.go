package types

import utils_fmt "exocomp/utils/fmt"
import _ "embed"
import net_url "net/url"
import "os"
import "strings"

type Config struct {
	Name        string
	Agent       string
	Model       string
	Playground  string
	Prompt      string
	Sandbox     string
	Temperature float64
	URL         *net_url.URL
}

func NewConfig(name string, agent string, model string, playground string, prompt string, sandbox string, temperature float64, url *net_url.URL) *Config {

	prompt = utils_fmt.FormatSingleLine(prompt)

	if temperature < 0.1 {
		temperature = 0.1
	} else if temperature > 1.0 {
		temperature = 1.0
	}

	if playground == "" {

		base := os.TempDir()
		tmp, err := os.MkdirTemp(base, "exocomp-playground-*")

		if err == nil {
			playground = tmp
		} else {
			playground = "/tmp/exocomp"
		}

	}

	return &Config{
		Name:        name,
		Agent:       agent,
		Model:       model,
		Playground:  playground,
		Prompt:      prompt,
		Sandbox:     sandbox,
		Temperature: temperature,
		URL:         url,
	}

}

func (config *Config) GetPrompt() string {
	return strings.TrimSpace(config.Prompt)
}

func (config *Config) ResolvePath(path string) *net_url.URL {

	endpoint := config.URL.ResolveReference(&net_url.URL{
		Path: path,
	})

	return endpoint

}
