package types

import "errors"
import "flag"
import net_url "net/url"

type Config struct {
	Model   string
	URL     *net_url.URL
	Sandbox string
	Verbose bool
}

func ParseConfig() (*Config, error) {

	tmp_model := flag.String(
		"model",
		"qwen3-coder:30b",
		"Ollama Model to use",
	)

	tmp_url := flag.String(
		"url",
		"http://localhost:11434/api",
		"Ollama API endpoint",
	)

	tmp_sandbox := flag.String(
		"sandbox",
		"/tmp/gomcp-project",
		"Project sandbox folder",
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

			return &Config{
				Model:   *tmp_model,
				URL:     url,
				Sandbox: *tmp_sandbox,
				Verbose: *tmp_verbose,
			}, nil

		} else {
			return nil, errors.New("Ollama URL must match (http|https)://.../api")
		}

	} else {
		return nil, err
	}

}

func (config *Config) ResolvePath(path string) *net_url.URL {

	endpoint := config.URL.ResolveReference(&net_url.URL{
		Path: path,
	})

	return endpoint

}
