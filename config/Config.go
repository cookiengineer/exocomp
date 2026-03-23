package config

import "exocomp/gadgets"
import _ "embed"
import "errors"
import "flag"
import net_url "net/url"
import "os"
import "slices"
import "strings"

//go:embed Config.prompt.txt
var config_prompt []byte

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

	cwd, _ := os.Getwd()

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
		cwd,
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

	prompts := make([]string, 0)

	files_help,    _ := gadgets.NewFiles(config.Sandbox).Help([]string{})
	programs_help, _ := gadgets.NewPrograms(config.Sandbox, config.Programs).Help([]string{})

	// TODO: Other Gadgets

	prompts = append(prompts, strings.TrimSpace(string(config_prompt)))
	prompts = append(prompts, "")
	prompts = append(prompts, "Available gadgets:")
	prompts = append(prompts, "")
	prompts = append(prompts, strings.TrimSpace(files_help))
	prompts = append(prompts, "")
	prompts = append(prompts, strings.TrimSpace(programs_help))

	return strings.Join(prompts, "\n")

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
