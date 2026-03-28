package config

import "exocomp/tools"
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

	// Exocomp settings
	Model       string
	URL         *net_url.URL
	Temperature float64
	Verbose     bool

	// Tool settings
	Agent       string
	Sandbox     string
	Tools       []string
	Programs    []string

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

	tmp_temperature := flag.Float64(
		"temperature",
		0.3,
		"0.3 code, 0.5 balanced, 0.7 creative, 1.0 hallucinations",
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

			allowed_tools    := make([]string, 0)
			allowed_programs := make([]string, 0)

			for _, tool := range AllowedTools {
				allowed_tools = append(allowed_tools, tool)
			}

			for _, program := range AllowedPrograms {
				allowed_programs = append(allowed_programs, program)
			}

			return &Config{

				Model:       *tmp_model,
				URL:         url,
				Temperature: temperature,
				Verbose:     *tmp_verbose,

				Agent:       *tmp_agent,
				Sandbox:     *tmp_sandbox,
				Tools:       allowed_tools,
				Programs:    allowed_programs,

			}, nil

		} else {
			return nil, errors.New("Ollama URL must match (http|https)://.../api")
		}

	} else {
		return nil, err
	}

}

func (config *Config) GetPrompt() string {

	// TODO: Integrate Tasks Help when ready

	help := make([]string, 0)

	for _, tool := range config.Tools {

		if tool == "files" {
			tmp, _ := tools.NewFiles(config.Agent, config.Sandbox, config.Tools, config.Programs).Help([]string{})
			help    = append(help, strings.TrimSpace(tmp))
		} else if tool == "notes" {
			tmp, _ := tools.NewNotes(config.Agent, config.Sandbox, config.Tools, config.Programs).Help([]string{})
			help    = append(help, strings.TrimSpace(tmp))
		} else if tool == "programs" {
			tmp, _ := tools.NewPrograms(config.Agent, config.Sandbox, config.Tools, config.Programs).Help([]string{})
			help    = append(help, strings.TrimSpace(tmp))
		} else if tool == "tasks" {
			// tmp, _ := tools.NewTasks(config.Agent, config.Sandbox, config.Tools, config.Programs).Help([]string{})
			// help    = append(help, strings.TrimSpace(tmp))
		}

	}

	return strings.Join([]string{
		strings.TrimSpace(string(config_prompt)),
		"",
		"You have access to the following tools:",
		"",
		strings.Join(help, "\n\n"),
	}, "\n")

}

func (config *Config) IsAllowedProgram(name string) bool {
	return slices.Contains(config.Programs, name)
}

func (config *Config) IsAllowedTool(name string) bool {
	return slices.Contains(config.Tools, name)
}

func (config *Config) ResolvePath(path string) *net_url.URL {

	endpoint := config.URL.ResolveReference(&net_url.URL{
		Path: path,
	})

	return endpoint

}
