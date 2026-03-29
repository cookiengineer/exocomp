package config

import _ "embed"
import "errors"
import "flag"
import net_url "net/url"
import "os"
import "strings"

//go:embed Config.prompt.txt
var config_prompt []byte

type Config struct {
	Agent       string
	Model       string
	Sandbox     string
	Temperature float64
	URL         *net_url.URL
	Verbose     bool
}

func ParseConfig() (*Config, error) {

	cwd, _ := os.Getwd()

	tmp_agent := flag.String(
		"agent",
		"manager",
		"Agent type: manager, coder, tester",
	)

	tmp_model := flag.String(
		"model",
		"qwen3-coder:30b",
		"Ollama Model to use",
	)

	tmp_sandbox := flag.String(
		"sandbox",
		cwd,
		"Project sandbox",
	)

	tmp_temperature := flag.Float64(
		"temperature",
		0.3,
		"0.3 code, 0.5 balanced, 0.7 creative, 1.0 hallucinations",
	)

	tmp_url := flag.String(
		"url",
		"http://localhost:11434/api",
		"Ollama API endpoint",
	)

	tmp_verbose := flag.Bool(
		"verbose",
		false,
		"Enable verbose logging",
	)

	flag.Parse()

	url, err := net_url.Parse(*tmp_url)

	if err == nil {

		if (url.Scheme == "http" || url.Scheme == "https") && url.Path == "/api" {

			temperature := *tmp_temperature

			if temperature < 0.1 {
				temperature = 0.1
			} else if temperature > 1.0 {
				temperature = 1.0
			}

			return &Config{
				Agent:       *tmp_agent,
				Model:       *tmp_model,
				Sandbox:     *tmp_sandbox,
				Temperature: temperature,
				URL:         url,
				Verbose:     *tmp_verbose,
			}, nil

		} else {
			return nil, errors.New("Ollama URL must match (http|https)://.../api")
		}

	} else {
		return nil, err
	}

}

func (config *Config) GetPrompt() string {

	return strings.Join([]string{
		strings.TrimSpace(string(config_prompt)),
	}, "\n")

}

func (config *Config) ResolvePath(path string) *net_url.URL {

	endpoint := config.URL.ResolveReference(&net_url.URL{
		Path: path,
	})

	return endpoint

}
