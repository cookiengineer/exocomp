package config

import _ "embed"
import "errors"
import "flag"
import net_url "net/url"
import "slices"

//go:embed Config.prompt.txt
var prompt []byte

type Config struct {
	Agent    string
	Model    string
	Gadgets  []string
	Programs []string
	URL      *net_url.URL
	Sandbox  string
	Verbose  bool
}

func ParseConfig() (*Config, error) {

	tmp_agent := flag.String(
		"agent",
		"coder",
		"Agent type: coder, tester, manager",
	)

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
		"/tmp/project-codebase",
		"Project sandbox",
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

			allowed_gadgets  := make([]string, 0)
			allowed_programs := AllowedPrograms[0:]

			for _, gadget_type := range AllowedGadgets {
				allowed_gadgets = append(allowed_gadgets, gadget_type.String())
			}

			return &Config{
				Agent:    *tmp_agent,
				Model:    *tmp_model,
				URL:      url,
				Gadgets:  allowed_gadgets,
				Programs: allowed_programs,
				Sandbox:  *tmp_sandbox,
				Verbose:  *tmp_verbose,
			}, nil

		} else {
			return nil, errors.New("Ollama URL must match (http|https)://.../api")
		}

	} else {
		return nil, err
	}

}

func (config *Config) GetPrompt() string {

	prompt := string(prompt)
	prompt += "The list of available gadgets is:"

	for _, name := range config.Gadgets {
		prompt += "#!" + name + "\n"
	}

	return prompt

}

func (config *Config) IsAllowedGadget(name string) bool {
	return slices.Contains(config.Gadgets, name)
}

func (config *Config) IsAllowedProgram(name string) bool {
	return slices.Contains(config.Programs, name)
}

func (config *Config) ResolvePath(path string) *net_url.URL {

	endpoint := config.URL.ResolveReference(&net_url.URL{
		Path: path,
	})

	return endpoint

}
